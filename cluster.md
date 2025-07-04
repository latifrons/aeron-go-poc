# Cluster Configs

## MediaDriverContext

| Key                | Value                                                         | Description                                                |
|--------------------|---------------------------------------------------------------|------------------------------------------------------------|
| aeronDirectoryName | C:\Users\LATFIR~1\AppData\Local\Temp\aeron-latfirons-0-driver | The directory where the media driver will store its files. |

## ReplicationArchiveContext

| Key                     | Value                          | Description |
|-------------------------|--------------------------------|-------------|
| controlResponseChannel  | aeron:udp?endpoint=localhost:0 |             |
| controlResponseStreamId | 20                             |             |

## ArchiveContext

| Key                  | Value                                                         | Description                                                |
|----------------------|---------------------------------------------------------------|------------------------------------------------------------|
| aeronDirectoryName   | C:\Users\LATFIR~1\AppData\Local\Temp\aeron-latfirons-0-driver | The directory where the media driver will store its files. |
| archiveDir           | D:\ws\latifrons\aeron-java-poc\aeron-cluster-0\archive        |                                                            |
| controlChannel       | aeron:udp?endpoint=localhost:10001                            |                                                            |
| archiveClientContext | [ReplicationArchiveContext]                                   |                                                            |
| localControlChannel  | aeron:ipc?term-length=64k                                     |                                                            |
| replicationChannel   | aeron:udp?endpoint=localhost:0                                |                                                            |

## AeronArchiveContext

| Key                    | Value                                 | Description |
|------------------------|---------------------------------------|-------------|
| lock                   | [NoOpLock]                            |             |
| controlRequestChannel  | [ArchiveContext.localControlChannel]  |             |
| controlRequestStreamId | [ArchiveContext.localControlStreamId] |             |
| controlRequestChannel  | [ArchiveContext.localControlChannel]  |             |
| controlRequestStreamId | [ArchiveContext.localControlStreamId] |             |
| archiveClientContext   | [ReplicationArchiveContext]           |             |
| localControlChannel    | aeron:ipc?term-length=64k             |             |
| replicationChannel     | aeron:udp?endpoint=localhost:0        |             |

## ConsensusModuleContext

