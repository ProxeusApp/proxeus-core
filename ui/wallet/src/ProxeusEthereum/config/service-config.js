const serviceConfig = {
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
