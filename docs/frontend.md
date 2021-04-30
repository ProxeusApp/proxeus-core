# Frontend Development

### Prerequisites
+ yarn (1.12+)
+ node (14+)
+ vue-cli

> **Important**
>
>  Only use yarn 1.12 or higher. For linking together local dependencies we use Yarn Workspaces:
> https://yarnpkg.com/lang/en/docs/workspaces/
>
> Use yarn only from the /core/central/ui directory, as the common dependencies will be stored in
> /core/central/ui/node_modules instead of the package subfolders (./core, etc.).

#### Typescript
Use typescript wherever you find it appropriate. Use it especially for critical
modules and components where static type checking is helpful.
Check https://vuejs.org/v2/guide/typescript.html for more information in regards to TypeScript in a Vue context.

#### Frontend Dev Server
~~~
yarn run serve
~~~
Webpack-dev-server runs on port 3005 which points to the local
frontend dev server, which points at the go backend.

~~~
##### BrowserSync (Port 3005) #####
                |
                |
                V
##### go webserver (Port 1323) #####
~~~

Important: If you work on :1323, the frontend code will reflect
the latest production build. Use 3005 if you want to use the latest frontend features


### Technologies used
- Vue https://vuejs.org/v2/guide/index.html
- Vue Router https://router.vuejs.org/
- Vuex https://github.com/vuejs/vuex
- axios https://github.com/axios/axios

- Webpack https://webpack.js.org/

- Progressive Web Apps https://developers.google.com/web/progressive-web-apps/
- FileSystem API https://developer.mozilla.org/en-US/docs/Web/API/FileSystem
- Web3 https://github.com/ethereum/web3.js/tree/v0.20.6
- ES6/ES7… https://codeburst.io/javascript-wtf-is-es6-es8-es-2017-ecmascript-dca859e4821c
- Async/Await https://medium.com/@rafaelvidaurre/truly-understanding-async-await-491dd580500e
- JavaScript https://developer.mozilla.org/bm/docs/Web/JavaScript

- Bootstrap https://getbootstrap.com/docs/4.1/getting-started/introduction/
