# Front End Build

Please read the [Build All](build_all.md) section before continuing.

Building Proxeus's web client requires Node.js & yarn.

#### Building for Production
```
make ui
```

#### Building for Development (with hot-reloading)
```
make ui-dev
```

#### Building for Development (as preview)
```
make ui-serve
```

You should now be able to access the web page on the port shown in the console.

If you want to check the frontend code navigate to the `ui` folder.

The backend will include the frontend files directly in the source code using `bindata`.

Follow [Back End Build](build_backend.md) for how to build the server.
