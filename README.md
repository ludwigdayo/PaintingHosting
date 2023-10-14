<p align="center">
<img src="./docs/images/PaintingHosting.svg" width=50%/ alt=""></p>



## Introduction

#### This project is used for image storage and display on the web

![introduce-image](docs/images/introduce-image.png)



## Features

- Upload images in browser
- Display images on the web
- Ability to comment on images
- Ability to add stars to pictures
- Classify and display by tags
- No login required and no traces left



## Build Setup

Building the project requires [**go**](https://go.dev/doc/install) and [**node**](https://nodejs.org/en/download/current) packages, Since it depends on [cgo](https://go.dev/blog/cgo) you also need to install [**gcc**](https://gcc.gnu.org/install/binaries.html) package

``` shell
# In the project root directory
go build .

# Switch working path
cd front-end-code

# install dependencies
npm install

# build for production with minification
npm run build
```



## How to deploy

- Create these directory as:

  ````
  ```
  Workspace
  ├── db
  ├── dist
  └── ImagesData
  ```
  ````

- Place the front-end build product in the *dist* folder

- Make sure the three folders are on the same as the executable file

  ```
  Workspace
  ├── db
  │   └── imageshow-server-database.db
  ├── dist
  │   ├── index.html
  │   └── static
  ├── ImagesData
  └── PaintingHosting-linux-amd64
  ```

- Now you can run the [executable](https://github.com/succerseng/PaintingHosting/releases/latest)

  ```shell
  $ ./PaintingHosting-linux-amd64 -p 80 -h 0.0.0.0
  ```

  

## Licence

[GPL-3.0](https://github.com/succerseng/PaintingHosting/blob/main/LICENSE)
