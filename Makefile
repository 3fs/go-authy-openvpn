AUTH_SCRIPT_SOURCE=https://github.com/matevzmihalic/auth-script-openvpn/archive/params.zip

INSTDIR=/usr/lib/authy
BUILD_DIR= build

all: $(BUILD_DIR)/go-authy-openvpn $(BUILD_DIR)/auth_script.so

auth-script-openvpn:
	mkdir auth-script-openvpn && wget -qO- $(AUTH_SCRIPT_SOURCE) | bsdtar -xvf- --strip-components 1 -C auth-script-openvpn

$(BUILD_DIR)/auth_script.so: auth-script-openvpn
	make -C auth-script-openvpn && mv auth-script-openvpn/auth_script.so $(BUILD_DIR)

$(BUILD_DIR)/go-authy-openvpn:
	go build -o $(BUILD_DIR)/go-authy-openvpn ./src

clean:
	rm -rf $(BUILD_DIR)

install: all
	mkdir -p $(DESTDIR)$(INSTDIR)
	cp $(BUILD_DIR)/* $(DESTDIR)$(INSTDIR)
	chmod 755 $(DESTDIR)$(INSTDIR)/*
	mkdir -p $(DESTDIR)/usr/sbin
	cp scripts/authy-vpn-add-user $(DESTDIR)/usr/sbin/authy-vpn-add-user
	chmod 700 $(DESTDIR)/usr/sbin/authy-vpn-add-user
	./scripts/post-install
