## Azure Event Hub Publisher

## usage 

```shell
./eventmaker eventhub --device <your device name>
                      --file conf/example.yaml \
                      --connect <your Event Hub entity connection string>
```

You can locate the Event Hub connection string using the Azure CLI 

```shell
az eventhubs eventhub authorization-rule keys list 
    --namespace-name <your Event Hub namespace name> \
    --eventhub-name <your Event Hub name> \
    --name <your Event Hub authorization rule name>
```


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](../LICENSE)


