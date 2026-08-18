#!/bin/sh

TEST_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$TEST_DIR/../../.." && pwd)
. "$TEST_DIR/testlib.sh"

CTL_SCRIPT="$REPO_ROOT/contrib/openwrt/files/usr/sbin/steamhostsyncctl"

setup_control() {
    new_sandbox
    install_mock "$TEST_DIR/mocks/uci" uci
    install_mock "$TEST_DIR/mocks/dnsmasq-service" dnsmasq-service
    install_mock "$TEST_DIR/mocks/init-service" init-service
    install_mock "$TEST_DIR/mocks/cron-service" cron-service

    export STEAMHOSTSYNC_STATE_DIR="$TEST_ROOT/var/lib/steamhostsync"
    export STEAMHOSTSYNC_UCI_BIN="$TEST_ROOT/bin/uci"
    export STEAMHOSTSYNC_DNSMASQ_BIN="$TEST_ROOT/bin/dnsmasq-service"
    export STEAMHOSTSYNC_INIT_BIN="$TEST_ROOT/bin/init-service"
    export STEAMHOSTSYNC_CRON_BIN="$TEST_ROOT/bin/cron-service"
    export STEAMHOSTSYNC_CRONTAB="$TEST_ROOT/etc/crontabs/root"
    export STEAMHOSTSYNC_CONTROL_LOCK_DIR="$TEST_ROOT/var/lock/steamhostsync-control.lock"
    export STEAMHOSTSYNC_UPDATE_BIN="$TEST_ROOT/bin/update"
    export MOCK_UCI_STATE="$TEST_ROOT/uci-addnhosts"
    export MOCK_DNSMASQ_COUNT="$TEST_ROOT/dnsmasq.count"
    export MOCK_INIT_CALLS="$TEST_ROOT/init.calls"
    export MOCK_CRON_CALLS="$TEST_ROOT/cron.calls"
    : >"$STEAMHOSTSYNC_CRONTAB"
}

test_setup_is_idempotent() {
    setup_control

    sh "$CTL_SCRIPT" setup || return 1
    sh "$CTL_SCRIPT" setup || return 1

    assert_line_count 1 "$STEAMHOSTSYNC_STATE_DIR/Hosts" "$MOCK_UCI_STATE" || return 1
    assert_line_count 1 '# steamhostsync managed task' "$STEAMHOSTSYNC_CRONTAB" || return 1
    assert_line_count 1 "$STEAMHOSTSYNC_UPDATE_BIN" "$STEAMHOSTSYNC_CRONTAB" || return 1
    assert_line_count 1 'enable' "$MOCK_CRON_CALLS" || return 1
    assert_line_count 1 'restart' "$MOCK_CRON_CALLS"
}

test_disable_removes_only_managed_entries() {
    setup_control
    printf '%s\n' '/tmp/other-hosts' >"$MOCK_UCI_STATE"
    printf '%s\n' '5 1 * * * /usr/bin/custom # steamhostsync managed task' >"$STEAMHOSTSYNC_CRONTAB"
    sh "$CTL_SCRIPT" setup || return 1

    sh "$CTL_SCRIPT" disable || return 1

    assert_line_count 1 '/tmp/other-hosts' "$MOCK_UCI_STATE" || return 1
    assert_line_count 0 "$STEAMHOSTSYNC_STATE_DIR/Hosts" "$MOCK_UCI_STATE" || return 1
    assert_line_count 1 '/usr/bin/custom' "$STEAMHOSTSYNC_CRONTAB" || return 1
    assert_line_count 0 "$STEAMHOSTSYNC_UPDATE_BIN" "$STEAMHOSTSYNC_CRONTAB"
}

test_control_lock_prevents_concurrent_setup() {
    setup_control
    mkdir "$STEAMHOSTSYNC_CONTROL_LOCK_DIR"
    printf '%s\n' "$$" >"$STEAMHOSTSYNC_CONTROL_LOCK_DIR/pid"

    sh "$CTL_SCRIPT" setup
    status=$?

    assert_eq 75 "$status" || return 1
    assert_file_absent "$MOCK_UCI_STATE" || return 1
    assert_line_count 0 "$STEAMHOSTSYNC_UPDATE_BIN" "$STEAMHOSTSYNC_CRONTAB"
}

test_package_scripts_never_reference_system_hosts() {
    scripts="$REPO_ROOT/contrib/openwrt/files/usr/libexec/steamhostsync-update $REPO_ROOT/contrib/openwrt/files/usr/sbin/steamhostsyncctl $REPO_ROOT/contrib/openwrt/files/etc/init.d/steamhostsync"
    for script in $scripts; do
        [ -f "$script" ] || fail "missing package script: $script" || return 1
        if grep -F '/etc/hosts' "$script" >/dev/null 2>&1; then
            fail "package script references /etc/hosts: $script"
            return 1
        fi
    done
}

run_test test_setup_is_idempotent
run_test test_disable_removes_only_managed_entries
run_test test_control_lock_prevents_concurrent_setup
run_test test_package_scripts_never_reference_system_hosts
finish_tests
