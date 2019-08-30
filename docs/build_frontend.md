# Front End Build

Please read the [Build All](build_all.md) section before continuing.

Building Proxeus's web client requires Node.js 8 & yarn 1.12.

#### Building for Production
```
make ui
```
#### Building for Development (with hot-reloading)
```
make ui-dev
```

You should now be able to access the web page on `http://localhost:3005`

If you want to check the frontend code navigate to the `ui` folder.

The backend will include the frontend files directly in the source code using `bindata`.

Follow [Back End Build](build_backend.md) for how to build the server.
