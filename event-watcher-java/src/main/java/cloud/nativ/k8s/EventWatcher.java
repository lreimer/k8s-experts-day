package cloud.nativ.k8s;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

import com.google.gson.reflect.TypeToken;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.kubernetes.client.openapi.ApiClient;
import io.kubernetes.client.openapi.ApiException;
import io.kubernetes.client.openapi.Configuration;
import io.kubernetes.client.openapi.apis.CoreV1Api;
import io.kubernetes.client.openapi.models.V1Namespace;
import io.kubernetes.client.openapi.models.V1Pod;
import io.kubernetes.client.openapi.models.V1PodList;
import io.kubernetes.client.util.Config;
import io.kubernetes.client.util.Watch;
import okhttp3.OkHttpClient;

public class EventWatcher {
    private static final Logger LOGGER = LoggerFactory.getLogger(EventWatcher.class);

    public static void main(String[] args) throws IOException, ApiException {
        ApiClient client = Config.defaultClient();
        OkHttpClient httpClient = client.getHttpClient().newBuilder().readTimeout(0, TimeUnit.SECONDS).build();
        client.setHttpClient(httpClient);
        Configuration.setDefaultApiClient(client);

        CoreV1Api api = new CoreV1Api();

        LOGGER.info("List of Pods in namespace 'default'");

        V1PodList list = api.listNamespacedPod("default", null, null, null, null, null, null, null, null, null);
        for (V1Pod item : list.getItems()) {
            LOGGER.info("* {}", item.getMetadata().getName());
        }

        LOGGER.info("");
        LOGGER.info("Watching for Pod events in namespace 'default'");

        Watch<V1Pod> watch = Watch.createWatch(client,
                api.listNamespacedPodCall("default", null, null, null, null, null, 5, null, null, Boolean.TRUE, null),
                new TypeToken<Watch.Response<V1Pod>>() {
                }.getType());

        try {
            for (Watch.Response<V1Pod> item : watch) {
                LOGGER.info("{} : {}", item.type, item.object.getMetadata().getName());
            }
        } finally {
            watch.close();
        }
    }
}
