const serviceConfig = {
  sepolia: {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0x61a26381a8ca72870ab4E4108d5D3982a89D7fd0',
    PROXEUS_FS_ADDRESS: '0x9bc518Fd070BE3DBB07399261688015744a7FB02',
    // When the PROXEUS_FS_ADDRESS changes add the old address to the list below so that file validation check also past contracts
    PROXEUS_FS_PAST_ADDRESSES: [
    ]
  },
  goerli: {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0x15FeA089CC48B4f4596242c138156e3B53579B37',
    PROXEUS_FS_ADDRESS: '0x66FF4FBF80D4a3C85a54974446309a2858221689',
    // When the PROXEUS_FS_ADDRESS changes add the old address to the list below so that file validation check also past contracts
    PROXEUS_FS_PAST_ADDRESSES: [
    ]
  },
  main: {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0xa017ac5fac5941f95010b12570b812c974469c2c',
    PROXEUS_FS_ADDRESS: '0xf63e471d8cbc57517c37c39c35381a385628e012',
    // When the PROXEUS_FS_ADDRESS changes add the old address to the list below so that file validation check also past contracts
    PROXEUS_FS_PAST_ADDRESSES: [
    ]
  },
  polygon: {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0x6B586cdC3f889AD4C9FA78000F99C54e88F66459',
    PROXEUS_FS_ADDRESS: '0x60970BeFda93464A105DD21Dc6a30B69C5B5c6e4',
    // When the PROXEUS_FS_ADDRESS changes add the old address to the list below so that file validation check also past contracts
    PROXEUS_FS_PAST_ADDRESSES: [
    ]
  },
  'polygon-mumbai': {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0xf94BdC648A30719fF0cf91B436f9819F7804e1a0',
    PROXEUS_FS_ADDRESS: '0x00119d8C02bbC4c1231D054BB2813792B4411Ed5',
    // When the PROXEUS_FS_ADDRESS changes add the old address to the list below so that file validation check also past contracts
    PROXEUS_FS_PAST_ADDRESSES: [
    ]
  }
}

export default serviceConfig
