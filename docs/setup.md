# Setup Instructions
Welcome! To get started with our project, you need to install a couple of essential tools:

- Go (version 1.22 or higher)
- Rye (a Python package manager)

Follow the instructions below to set up your environment.

### Go

You will find the instruction to install Go depending on your infrastructure here:
https://go.dev/doc/install

### Rye

You will find the instruction to install Rye depending on your infrastructure here:
https://rye.astral.sh/guide/installation/

### Rust

You will find the instruction to install Rust here:
https://www.rust-lang.org/tools/install

### Setup environment variables

In the `/templates` folder you will find the templates of the environment variables
needed to run the project. If you do not want to modify it, you can just run:

```bash
cp ./templates/.env.template ./.env
cp ./templates/grafana.env.template ./grafana.env
cp ./templates/production.env.template ./production.env
```
You will need the environment variables defined in the `.env` file to be loader.
You can either use the [dotenv] plugin or load it using:

```bash
export $(grep -v '^#' .env | xargs)
```
### Next steps

At the root of the repository, you will find a `makefile` that assists with
various tasks.

#### Download the data

The first step is to download the data that you will find [here](https://www.kaggle.com/datasets/nphantawee/pump-sensor-data).
Put it inside the `/data` directory.

#### Build the Model

The second step is to train the autoencoder. To do this, go to the `ml_model`
directory and run:
```bash
make model/train
```
This command will train the model and save the artifacts in the `ml_service`
directory.

#### Build the Facade

Next, you'll need to containerize the ML facade and the ML service. To do this, run:
```bash
make services/build
```

#### Run the Services

Once everything is trained and built, you can start the services by running:
```bash
make services/run
```

To stop the services, run:
```bash
make services/stop
```

To delete the containers, run:
```bash
make services/down
```

#### Send the data

You can send data to either the RabbitMQ server or the API endpoint by running:
```bash
make run/send-data --rabbitmq={true/false} --requests=10
```

#### Check the dashboards

You can log in to [Grafana](http://localhost:9000) to check the dashboards and monitor predictions in real-time.

<!--references-->
[dotenv]: https://github.com/ohmyzsh/ohmyzsh/blob/master/plugins/dotenv/dotenv.plugin.zsh
