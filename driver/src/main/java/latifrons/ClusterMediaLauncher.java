package latifrons;

import io.aeron.cluster.ClusteredMediaDriver;
import io.aeron.driver.ThreadingMode;
import org.agrona.ErrorHandler;
import org.agrona.concurrent.BusySpinIdleStrategy;
import org.agrona.concurrent.NoOpIdleStrategy;
import org.agrona.concurrent.ShutdownSignalBarrier;

import java.io.File;
import java.util.Arrays;
import java.util.List;

import static java.lang.Integer.parseInt;

public class ClusterMediaLauncher {
    private static ErrorHandler errorHandler(final String context) {
        return
                (Throwable throwable) ->
                {
                    System.err.println(context);
                    throwable.printStackTrace(System.err);
                };
    }

    private static final int PORT_BASE = 9000;

    public static void launch() {
        final int nodeId = parseInt(EnvTool.mustGetEnv("aeron.cluster.nodeId"));
        final String hostnamesStr = EnvTool.mustGetEnv("aeron.cluster.hostnames");
        final String clusterDir = EnvTool.mustGetEnv("aeron.cluster.dir");
        final String aeronDir = EnvTool.mustGetEnv("aeron.driver.dir");
        final String aeronIdle = EnvTool.getEnv("aeron.driver.idle", "");
        final String lowLatency = EnvTool.getEnv("aeron.driver.lowLatency", "0");

        final int ingressStreamId = parseInt(EnvTool.mustGetEnv("aeron.cluster.ingressStreamId"));
        final List<String> hostnames = Arrays.asList(hostnamesStr.split(","));
        final List<String> internalHostnames = Arrays.asList(hostnamesStr.split(","));

        final ShutdownSignalBarrier barrier = new ShutdownSignalBarrier();

        final ClusterConfig clusterConfig = ClusterConfig.create(
                nodeId, hostnames, internalHostnames,
                PORT_BASE,
                null);

        clusterConfig.mediaDriverContext().
                aeronDirectoryName(aeronDir).
                errorHandler(errorHandler("Media Driver"));

        if (lowLatency.equals("1")) {
            System.out.println("Using low latency settings for Aeron Media Driver.");
            clusterConfig.mediaDriverContext()
                    .threadingMode(ThreadingMode.DEDICATED)
                    .termBufferSparseFile(false)
                    .socketRcvbufLength(2 * 1024 * 1024)
                    .socketSndbufLength(2 * 1024 * 1024)
                    .initialWindowLength(2 * 1024 * 1024)
                    .conductorIdleStrategy(new BusySpinIdleStrategy())
                    .senderIdleStrategy(new NoOpIdleStrategy())
                    .receiverIdleStrategy(new NoOpIdleStrategy());
        } else {
            System.out.println("Using default settings for Aeron Media Driver.");
            clusterConfig.mediaDriverContext().
                    senderIdleStrategy(ClusterConfig.toIdleStrategy(aeronIdle)).
                    receiverIdleStrategy(ClusterConfig.toIdleStrategy(aeronIdle)).
                    conductorIdleStrategy(ClusterConfig.toIdleStrategy(aeronIdle));
        }

        clusterConfig.archiveContext()
                .archiveDir(new File(clusterDir, "archive"))
                .errorHandler(errorHandler("Archive"));
        clusterConfig.aeronArchiveContext()
                .errorHandler(errorHandler("Aeron Archive"));
        clusterConfig.consensusModuleContext()
                .clusterDir(new File(clusterDir, "cluster"))
                .ingressChannel(ClusterConfig.udpChannel(nodeId, hostnames.get(nodeId), PORT_BASE, ClusterConfig.CLIENT_FACING_PORT_OFFSET))
                .ingressStreamId(ingressStreamId)
                .errorHandler(errorHandler("Consensus Module"));
        clusterConfig.clusteredServiceContext()
                .errorHandler(errorHandler("Clustered Service"));

        try (
                ClusteredMediaDriver ignore = ClusteredMediaDriver.launch(
                        clusterConfig.mediaDriverContext(),
                        clusterConfig.archiveContext(),
                        clusterConfig.consensusModuleContext())) {

            // ClusteredServiceContainer ignore2 = ClusteredServiceContainer.launch(
            //                        clusterConfig.clusteredServiceContext())
            System.out.println("[" + nodeId + "] Started Cluster Node on " + hostnames.get(nodeId) + "...");
            System.out.println("[" + nodeId + "] Cluster Members: " + clusterConfig.consensusModuleContext().clusterMembers());
            System.out.println("[" + nodeId + "] Aeron Folder: " + clusterConfig.mediaDriverContext().aeronDirectoryName());
            System.out.println("[" + nodeId + "] Cluster Folder: " + clusterConfig.consensusModuleContext().clusterDir());
            System.out.println("[" + nodeId + "] Archive Folder: " + clusterConfig.archiveContext().archiveDirectoryName());

            barrier.await();
            System.out.println("[" + nodeId + "] Exiting");
        } catch (Exception e) {
            System.err.println("[" + nodeId + "] Error during Cluster Node startup: " + e.getMessage());
            e.printStackTrace(System.err);
            throw e;
        }
    }
}
