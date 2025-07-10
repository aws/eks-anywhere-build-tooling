/* Custom iPXE configuration for EKS Anywhere bare metal provisioning */

#define LINK_WAIT_MS 15000    /* Increased from 5000ms for better stability */
#define NET_PROTO_IPV4
#define DOWNLOAD_PROTO_TFTP
#define DOWNLOAD_PROTO_HTTP
#define IMAGE_EFI
#define CONSOLE_EFI
#define CONSOLE_SERIAL

#define CERT_CMD              /* Certificate management commands */
#define DIGEST_CMD            /* Image crypto digest commands */
#define DOWNLOAD_PROTO_HTTPS  /* Secure Hypertext Transfer Protocol */
#define IMAGE_TRUST_CMD       /* Image trust management commands */
#define NET_PROTO_IPV6        /* IPv6 protocol */
#define NSLOOKUP_CMD          /* DNS resolving command */
#define NTP_CMD               /* NTP commands */
#define NVO_CMD               /* Non-volatile option storage commands */
#define PARAM_CMD             /* params and param commands, for POSTing to tink */
#define PING_CMD              /* Ping command */
#define POWEROFF_CMD          /* Power off command */
#define REBOOT_CMD            /* Reboot command */
#define SANBOOT_PROTO_HTTP    /* HTTP SAN protocol */
#define VLAN_CMD              /* VLAN commands */
#define DOWNLOAD_PROTO_NFS    /* NFS */
#define ROUTE_CMD             /* Routing table management commands */
#define NET_PROTO_LACP        /* Link Aggregation control protocol */


/* Keep MAX_MODULES from previous configuration */
#define MAX_MODULES 17

/* Explicitly disable features based on upstream */
#undef SANBOOT_PROTO_AOE      /* AoE protocol */
#undef SANBOOT_PROTO_FCP      /* Fibre Channel protocol */
#undef SANBOOT_PROTO_IB_SRP   /* Infiniband SCSI RDMA protocol */
#undef SANBOOT_PROTO_ISCSI    /* iSCSI protocol */
#undef USB_EFI                /* Provide EFI_USB_IO_PROTOCOL interface */
#undef USB_HCD_EHCI           /* EHCI USB host controller */
#undef USB_HCD_UHCI           /* UHCI USB host controller */
#undef USB_HCD_XHCI           /* xHCI USB host controller */
#undef USB_KEYBOARD           /* USB keyboards */
#undef NET_PROTO_EAPOL        /* Workaround for Mellanox issue */
