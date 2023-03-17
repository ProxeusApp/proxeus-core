const serviceConfig = {
  goerli: {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0x15FeA089CC48B4f4596242c138156e3B53579B37',
    PROXEUS_FS_ADDRESS: '0x66FF4FBF80D4a3C85a54974446309a2858221689',
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
  ropsten: {
    DEFAULT_GAS_REGULAR: 1000000,
    DEFAULT_GAS_CREATE: 4000000,
    XES_TOKEN_ADDRESS: '0x84E0b37e8f5B4B86d5d299b0B0e33686405A3919',
    PROXEUS_FS_ADDRESS: '0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71', // Eternal Storage at: 0x7b83acb323fd4bbd874c1e9c295e0f486d94b233
    // When the PROXEUS_FS_ADDRESS changes add the old address to the list below so that file validation check also past contracts
    PROXEUS_FS_PAST_ADDRESSES: [
    ]
  }
}

export default serviceConfig
