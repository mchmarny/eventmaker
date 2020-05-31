## GCP PubSub Publisher

## usage 

> Note, the path to your authentication credentials must be defined in GOOGLE_APPLICATION_CREDENTIALS environment variable. This is the file path of your JSON file that contains your service account key. 

```shell
./eventmaker eventhub --device <your device name>
                      --file conf/example.yaml \
                      --project <your GCP project ID> \
                      --topic <your topic name in the above project>
```

For more information on how to configure your GCP project authentication see [here](https://cloud.google.com/docs/authentication/getting-started)

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](../LICENSE)


