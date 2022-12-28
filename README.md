# go-caching-with-redis

An example app I wrote while exploring caching with redis

## Setup

### Start Docker Containers

-   This will start `redis`, `postgres` and `adminer` containers

    ```sh
    docker compose up
    ```

### Download and Load sample dataset

(There could be other better/simpler ways to do this)

-   Download sample dataset (for more datasets check [IMDb Datasets][imdb_datasets])

    ```sh
    mkdir ./datasets/
    wget -O ./datasets/title.basics.tsv.gz https://datasets.imdbws.com/title.basics.tsv.gz
    ```

-   unzip & get the top 5000 rows (skip the header)

    ```sh
    gunzip --stdout ./datasets/title.basics.tsv.gz > ./datasets/title.basics.tsv \
        && head -5001 ./datasets/title.basics.tsv | tail -5000 > ./datasets/title_basics_top5k.tsv
    ```

    ```sh
    # Add ./datasets directory to .gitignore
    echo "\ndatasets/*" >> ./.gitignore

    # **Optional**: remove original files
    rm -ir ./datasets/title.basics.tsv ./datasets/title.basics.tsv.gz
    ```

-   Run the `CREATE TABLE` SQL command in `Adminer`. (For default docker container credentials check the .env.template file)

    ```SQL
    CREATE TABLE IF NOT EXISTS title_basics (
        tconst VARCHAR(10) PRIMARY KEY UNIQUE NOT NULL,
        title_type VARCHAR(20) NOT NULL,
        primary_title VARCHAR(100) NOT NULL,
        original_title VARCHAR(100) NOT NULL,
        is_adult VARCHAR(1) NOT NULL,
        start_year VARCHAR(4),
        end_year VARCHAR(4),
        runtime_minutes VARCHAR(5),
        genres VARCHAR(100)
    )
    ```

-   Import the data into the database using Adminer -> Click the table name -> Select data -> Import (as `TSV`)

## Teardown

### Stop Docker Containers

-   This will stop redis, postgres and adminer containers

    ```sh
    docker compose down
    ```

## API

-   Get Title information

    ```sh
    curl -X GET localhost:3000/title/tt0000001 \
        -H 'Content-Type: application/json' | jq .
    ```

[imdb_datasets]: https://www.imdb.com/interfaces/
