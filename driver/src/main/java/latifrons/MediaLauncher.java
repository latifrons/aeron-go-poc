package latifrons;

import io.aeron.Aeron;
import io.aeron.driver.MediaDriver;
import io.aeron.driver.ThreadingMode;
import org.agrona.ErrorHandler;
import org.agrona.concurrent.BusySpinIdleStrategy;
import org.agrona.concurrent.NoOpIdleStrategy;
import org.agrona.concurrent.ShutdownSignalBarrier;

import java.io.File;

public class MediaLauncher {
    private static ErrorHandler errorHandler(final String context) {
        return
                (Throwable throwable) ->
                {
                    System.err.println(context);
                    throwable.printStackTrace(System.err);
                };
    }

    public static void launch() {
        final String aeronDir = EnvTool.mustGetEnv("aeron.aeronDir");
        final String aeronIdle = EnvTool.getEnv("aeron.driver.idle", "");
        final String lowLatency = EnvTool.getEnv("aeron.driver.lowLatency", "0");

        MediaDriver.Context mediaDriverCtx = new MediaDriver.Context()
                .aeronDirectoryName(aeronDir)
                .dirDeleteOnStart(true)
                .dirDeleteOnShutdown(true)
                .errorHandler(errorHandler("Media Driver Error"));

        if (lowLatency.equals("1")) {
            System.out.println("Using low latency settings for Aeron Media Driver.");
            //[done] aeron.term.buffer.sparse.file=false
            //[client] aeron.pre.touch.mapped.memory=true
            //[done] aeron.socket.so_sndbuf=2m
            //[done] aeron.socket.so_rcvbuf=2m
            //[done] aeron.rcv.initial.window.length=2m
            //[done] aeron.threading.mode=DEDICATED
            //[done] aeron.sender.idle.strategy=noop
            //[done] aeron.receiver.idle.strategy=noop
            //[done] aeron.conductor.idle.strategy=spin
            //agrona.disable.bounds.checks=true

            mediaDriverCtx = mediaDriverCtx
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

            mediaDriverCtx = mediaDriverCtx
                    .senderIdleStrategy(ClusterConfig.toIdleStrategy(aeronIdle))
                    .receiverIdleStrategy(ClusterConfig.toIdleStrategy(aeronIdle))
                    .conductorIdleStrategy(ClusterConfig.toIdleStrategy(aeronIdle));
        }

        MediaDriver mediaDriver = MediaDriver.launch(mediaDriverCtx);
        // keep the media driver running

        System.out.println("Media driver started at " + new File(mediaDriver.aeronDirectoryName()).getAbsolutePath());
        System.out.println("Media Driver is running. Press Ctrl+C to exit.");
        try {
            Thread.currentThread().join();  // Keep the main thread alive
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.err.println("Media Driver interrupted: " + e.getMessage());
        }
        mediaDriver.close();
    }
}