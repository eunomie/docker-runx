# docker runx - hello example

This is a simple example of `docker runx` usage.

The main idea is to use an `alpine` image and create two actions:
- one, `user` that will get the user name from the environment variable
- the other, `ask` that will ask the user for its name

## Creation of the image

This `README.md` file will be used as the documentation for the image.

The [`runx.yaml`](runx.yaml) file contains the definition of the actions.

```
$ docker runx decorate alpine --tag NAMESPACE/REPOSITORY
```

## Usage

### Print this file

```
$ docker runx NAMESPACE/REPOSITORY --docs
```

### Run the `user` action

```
$ docker runx NAMESPACE/REPOSITORY user
```

### Run the `ask` action

```
$ docker runx NAMESPACE/REPOSITORY ask
```

You can also provide the user name as an argument:

```
$ docker runx NAMESPACE/REPOSITORY ask --opt name=John
```
