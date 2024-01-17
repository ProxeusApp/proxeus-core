/*

*** DYNAMIC CONFIG ***

Use this file to make dynamic configuration for the application.

How to use:

1) Pass VUE_APP_USE_DYNAMIC_CONFIG variable while build pipeline. Like this:
$ export VUE_APP_USE_DYNAMIC_CONFIG=true
2) Place this file into ui/core/dist/static/dynamic-config directory. You can mount Docker volume directry into this path. Any extra files can be placed as well.

Example of this configuration file:

export default {
  "apply": true, // Allow or deny this config (even with VUE_APP_USE_DYNAMIC_CONFIG value)
  "company": {
    "logo": {
      "lightTheme": {
        "path": "/static/dynamic-config/..."
      },
      "darkTheme": {
        "path": "/static/dynamic-config/..."
      },
      "style": {
        "home": {
          ... css styles ...
        },
        "adminDashboard": {
          "full": {
            ... css styles ...
          },
          "small": {
            ... css styles ...
          }
        },
        "userDashboard": {
          "full": {
            ... css styles ...
          },
          "small": {
            ... css styles ...
          }
        }
      }
    },
    "copyrightText": {
      "lang": {
        "en": "..."
        ... another languages ...
      },
      "logo": {
        "style": {
          ... css styles ...
        }
      },
      "link": "..."
    }
  },
  "homePage": {
    "platformTitle": {
      "lang": {
        "en": "..."
        ... another languages ...
      }
    },
    "platformDescription": {
      "lang": {
        "en": "..."
        ... another languages ...
      }
    }
  }
}

*/

export default {
  "apply": false
}
