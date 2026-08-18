include $(TOPDIR)/rules.mk

PKG_NAME:=steamhostsync-openwrt
PKG_VERSION:=0.1.0
PKG_RELEASE:=1
PKGARCH:=all

include $(INCLUDE_DIR)/package.mk

define Package/steamhostsync-openwrt
  SECTION:=net
  CATEGORY:=Network
  TITLE:=SteamHostSync OpenWrt updater
  URL:=https://github.com/Clov614/SteamHostSync
  DEPENDS:=+ca-bundle +uclient-fetch +dnsmasq
endef

define Package/steamhostsync-openwrt/description
  Architecture-neutral OpenWrt updater for SteamHostSync hosts subscriptions.
  Remote subscription is the default; local generation is optional.
endef

define Package/steamhostsync-openwrt/conffiles
/etc/config/steamhostsync
/etc/steamhostsync/allowed-domains
/etc/steamhostsync/router-config.yaml
endef

define Package/steamhostsync-openwrt/prerm
#!/bin/sh
[ -n "$${IPKG_INSTROOT}" ] || /usr/sbin/steamhostsyncctl disable || true
endef

define Package/steamhostsync-openwrt/install
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_DIR) $(1)/etc/steamhostsync
	$(INSTALL_DIR) $(1)/usr/libexec
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_CONF) ./files/etc/config/steamhostsync $(1)/etc/config/steamhostsync
	$(INSTALL_BIN) ./files/etc/init.d/steamhostsync $(1)/etc/init.d/steamhostsync
	$(INSTALL_BIN) ./files/usr/libexec/steamhostsync-update $(1)/usr/libexec/steamhostsync-update
	$(INSTALL_BIN) ./files/usr/sbin/steamhostsyncctl $(1)/usr/sbin/steamhostsyncctl
	$(INSTALL_CONF) ./files/etc/steamhostsync/allowed-domains $(1)/etc/steamhostsync/allowed-domains
	$(INSTALL_CONF) ./files/etc/steamhostsync/router-config.yaml $(1)/etc/steamhostsync/router-config.yaml
endef

$(eval $(call BuildPackage,steamhostsync-openwrt))
