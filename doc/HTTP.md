## HTTP Publisher 

## usage 

```shell
./eventmaker http --device <your device name>
                  --file conf/example.yaml \
                  --url <URL to POST to>
```

`eventmaker` will HTTP POST the generated event as `application/json` content. Anything else than a `200` HTTP status code response will be consider as an error


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](../LICENSE)


