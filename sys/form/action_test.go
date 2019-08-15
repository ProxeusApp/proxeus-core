package form

var complexActionFormSrc = `

{
  "inj0j46mbcukiw0xaxrvqv9tny4": {
    "name": "enterPat",
    "validate": {
      "required": true
    },
    "placeholder": "Enter Pat",
    "label": "Enter Pat",
    "help": "",
    "id": "bes6tsszphwxuxecpko6c",
    "action": {
      "source": [
        {
          "_destCompId": "qm3fnhob37k4c4affbiknj",
          "_index": 0,
          "comment": "",
          "regex": "Pat",
          "_fbonly_uipp": "4sjlas2mxkxxkg99maab0i"
        },
        {
          "_destCompId": "ui58usx36ijcbsrpzscn6s",
          "_index": 0,
          "comment": "",
          "regex": 52,
          "_fbonly_uipp": "4sjlas2mxkxxkg99maab0i"
        },
        {
          "_destCompId": "kmx39ke0vnmai4py4tuyp",
          "_index": 0,
          "comment": "",
          "regex": 53,
          "_fbonly_uipp": "4sjlas2mxkxxkg99maab0i"
        },
        {
          "_destCompId": "gd2ni8gwkip94mm5yzuhb",
          "_index": 0,
          "comment": "",
          "regex": "Antoine",
          "_fbonly_uipp": "4sjlas2mxkxxkg99maab0i"
        }
      ]
    },
    "_compId": "HC1",
    "_order": 2
  },
  "ayrk60w7t7sxtn10mu10p": {
    "name": "rootCheck",
    "validate": {
      "required": true
    },
    "help": "",
    "label": "Radio/Checkbox",
    "type": {
      "selected": 1,
      "all": [
        "radio",
        "checkbox"
      ]
    },
    "orientation": {
      "selected": 1,
      "all": [
        "horizontal",
        "vertical"
      ]
    },
    "values": [
      {
        "label": "Element 1",
        "value": 1,
        "help": ""
      },
      {
        "label": "Element 2",
        "value": 2,
        "help": ""
      },
      {
        "label": "Element 3",
        "value": 3,
        "help": ""
      },
      {
        "label": "Element 4",
        "value": 4,
        "help": ""
      },
      {
        "label": "Element 5",
        "value": 5,
        "help": ""
      }
    ],
    "_compId": "HC3",
    "_order": 1,
    "action": {
      "source": [
        {
          "_destCompId": "qm3fnhob37k4c4affbiknj",
          "_index": 0,
          "comment": "",
          "regex": "1",
          "_fbonly_uipp": "rfzf9kdr5gf5ny7y26l5t"
        },
        {
          "_destCompId": "fpu2k1ffyiuzvood2whqgf",
          "_index": 1,
          "comment": "",
          "regex": "2",
          "_fbonly_uipp": "bz4x9kzrq4cibz9da9hzx"
        },
        {
          "_destCompId": "f2i4iyj3rpmaaxhz3am6h",
          "_index": 2,
          "comment": "",
          "regex": "3",
          "_fbonly_uipp": "klgrostw7zpiic4ix79fgj"
        },
        {
          "_destCompId": "3vm7bzn2jnxalc4s6yga8k",
          "_index": 3,
          "comment": "",
          "regex": "4",
          "_fbonly_uipp": "21axjmmfsetq1g3933hsh"
        },
        {
          "_destCompId": "uincma9kwwle740ukhzt9t",
          "_index": 4,
          "comment": "",
          "regex": "5",
          "_fbonly_uipp": "qk3705g4wmx4dbo87vywn"
        }
      ]
    }
  },
  "qm3fnhob37k4c4affbiknj": {
    "name": "element1",
    "validate": {
      "required": true
    },
    "label": "Element 1",
    "values": [
      {
        "label": "Option 1",
        "value": "val1"
      },
      {
        "label": "Option 2",
        "value": "val2"
      }
    ],
    "action": {
      "destination": {
        "transition": {
          "selected": 1,
          "all": [
            "none",
            "slide",
            "fade"
          ]
        }
      },
      "source": [
        {
          "_destCompId": "fpu2k1ffyiuzvood2whqgf",
          "_index": 0,
          "comment": "",
          "regex": "val2",
          "_fbonly_uipp": "4abipc7v3la1xvii4vv4ot"
        }
      ]
    },
    "_compId": "HC4",
    "_order": 3
  },
  "fpu2k1ffyiuzvood2whqgf": {
    "name": "element2",
    "validate": {
      "required": true
    },
    "label": "Element 2",
    "values": [
      {
        "label": "Option 1",
        "value": "val1"
      },
      {
        "label": "Option 2",
        "value": "val2"
      }
    ],
    "action": {
      "destination": {
        "transition": {
          "selected": 1,
          "all": [
            "none",
            "slide",
            "fade"
          ]
        }
      },
      "source": [
        {
          "_destCompId": "f2i4iyj3rpmaaxhz3am6h",
          "_index": 0,
          "comment": "",
          "regex": "val2",
          "_fbonly_uipp": "q6ndpfwbz7mekdtku57las"
        }
      ]
    },
    "_compId": "HC4",
    "_order": 4
  },
  "f2i4iyj3rpmaaxhz3am6h": {
    "name": "element3",
    "validate": {
      "required": true
    },
    "label": "Element 3",
    "values": [
      {
        "label": "Option 1",
        "value": "val1"
      },
      {
        "label": "Option 2",
        "value": "val2"
      }
    ],
    "action": {
      "destination": {
        "transition": {
          "selected": 1,
          "all": [
            "none",
            "slide",
            "fade"
          ]
        }
      },
      "source": [
        {
          "_destCompId": "3vm7bzn2jnxalc4s6yga8k",
          "_index": 0,
          "comment": "",
          "regex": "val2",
          "_fbonly_uipp": "6wx3eg7ykfd6xo8s9ie1qd"
        }
      ]
    },
    "_compId": "HC4",
    "_order": 5
  },
  "3vm7bzn2jnxalc4s6yga8k": {
    "name": "element4",
    "validate": {
      "required": true
    },
    "label": "Element 4",
    "values": [
      {
        "label": "Option 1",
        "value": "val1"
      },
      {
        "label": "Option 2",
        "value": "val2"
      },
      {
        "label": "Option 3",
        "value": "val3"
      },
      {
        "label": "Option 4",
        "value": "val4"
      }
    ],
    "action": {
      "destination": {
        "transition": {
          "selected": 1,
          "all": [
            "none",
            "slide",
            "fade"
          ]
        }
      },
      "source": [
        {
          "_destCompId": "uincma9kwwle740ukhzt9t",
          "_index": 0,
          "comment": "",
          "regex": "val2",
          "_fbonly_uipp": "yxwqk6k96n44ll6sudn5h"
        }
      ]
    },
    "_compId": "HC4",
    "_order": 6
  },
  "uincma9kwwle740ukhzt9t": {
    "name": "element5",
    "validate": {
      "required": true
    },
    "label": "Element 5",
    "values": [
      {
        "label": "Option 1",
        "value": "val1"
      },
      {
        "label": "Option 2",
        "value": "val2"
      },
      {
        "label": "Option 3",
        "value": "val3"
      },
      {
        "label": "Option 4",
        "value": "val4"
      }
    ],
    "action": {
      "destination": {
        "transition": {
          "selected": 1,
          "all": [
            "none",
            "slide",
            "fade"
          ]
        }
      },
      "source": [
        {
          "_destCompId": "ui58usx36ijcbsrpzscn6s",
          "_index": 0,
          "comment": "",
          "regex": "val2",
          "_fbonly_uipp": "wsrr5jhkcvrxjrtpmhiwb"
        },
        {
          "_destCompId": "kmx39ke0vnmai4py4tuyp",
          "_index": 0,
          "comment": "",
          "regex": "val3",
          "_fbonly_uipp": "wsrr5jhkcvrxjrtpmhiwb"
        },
        {
          "_destCompId": "gd2ni8gwkip94mm5yzuhb",
          "_index": 0,
          "comment": "",
          "regex": "val4",
          "_fbonly_uipp": "wsrr5jhkcvrxjrtpmhiwb"
        }
      ]
    },
    "_compId": "HC4",
    "_order": 7
  },
  "ui58usx36ijcbsrpzscn6s": [
    {
      "name": "bla1",
      "validate": {
        "required": true
      },
      "label": "Element 5 Option 2",
      "placeholder": "Placeholder",
      "helps": [
        "describe all"
      ],
      "action": {
        "destination": {
          "transition": {
            "selected": 1,
            "all": [
              "none",
              "slide",
              "fade"
            ]
          }
        }
      },
      "_compId": "HC15",
      "_order": 8
    },
    {
      "name": "bla2",
      "validate": {
        "required": true
      },
      "label": "",
      "placeholder": "",
      "helps": [
        ""
      ],
      "_compId": "",
      "_order": ""
    },
    {
      "name": "bla3",
      "validate": {
        "required": true
      },
      "label": "",
      "placeholder": "",
      "helps": [
        ""
      ],
      "_compId": "",
      "_order": ""
    }
  ],
  "kmx39ke0vnmai4py4tuyp": [
    {
      "name": "bla31",
      "validate": {
        "required": true
      },
      "label": "Element 5 Option 3",
      "placeholder": "Placeholder",
      "helps": [
        "describe all"
      ],
      "action": {
        "destination": {
          "transition": {
            "selected": 1,
            "all": [
              "none",
              "slide",
              "fade"
            ]
          }
        }
      },
      "_compId": "HC15",
      "_order": 9
    },
    {
      "name": "bla32",
      "validate": {
        "required": true
      },
      "label": "",
      "placeholder": "",
      "helps": [
        ""
      ],
      "_compId": "",
      "_order": ""
    },
    {
      "name": "bla33",
      "validate": {
        "required": true
      },
      "label": "",
      "placeholder": "",
      "helps": [
        ""
      ],
      "_compId": "",
      "_order": ""
    }
  ],
  "gd2ni8gwkip94mm5yzuhb": [
    {
      "name": "bla41",
      "validate": {
        "required": true
      },
      "label": "Element 5 Option 4",
      "placeholder": "Placeholder",
      "helps": [
        "describe all"
      ],
      "action": {
        "destination": {
          "transition": {
            "selected": 1,
            "all": [
              "none",
              "slide",
              "fade"
            ]
          }
        }
      },
      "_compId": "HC15",
      "_order": 10
    },
    {
      "name": "bla42",
      "validate": {
        "required": true
      },
      "label": "",
      "placeholder": "",
      "helps": [
        ""
      ],
      "_compId": "",
      "_order": ""
    },
    {
      "name": "bla43",
      "validate": {
        "required": true
      },
      "label": "",
      "placeholder": "",
      "helps": [
        ""
      ],
      "_compId": "",
      "_order": ""
    }
  ],
  "8wayudwlw5sqlws6vxvvi": {
    "name": "phoneNrTest",
    "validate": {
      "required": true,
      "phoneNr": true
    },
    "placeholder": "Placeholder",
    "label": "Label",
    "help": "help text",
    "_compId": "HC1",
    "_order": 0
  }
}


`
