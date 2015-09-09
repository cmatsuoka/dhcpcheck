DHCPCheck
=========

DHCPCheck is an utility written to help network administrators to find
rogue DHCP servers and debug DHCP problems. DHCPCheck can run in many
modes, including **discover** (broadcast DHCPDISCOVER packets and show
DHCPOFFER packets and their originators), **snoop** (listen to any DHCP
packets running in the network), and **sentry** (warns the network
administrator if any anomalies are found, such as nonresponding servers
or unexpected responses coming from rogue servers).

This utility is currently under development.


Examples
--------

Send discover packet and show offers:
::

  # dhcpcheck discover -i wlp3s0 

