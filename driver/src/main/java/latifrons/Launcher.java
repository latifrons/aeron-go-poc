package latifrons;

public class Launcher {
    public static void main(String[] args) {
        if (args.length < 2) {
            System.out.println("Usage: java latifrons.Launcher <command> [<args>]");
            System.out.println("Available commands:");
            System.out.println("  solo - Start the solo server");
            System.out.println("  cluster - Start the cluster server");
            return;
        }
        String command = args[1];

        switch (command) {
            case "solo":
                MediaLauncher.launch();
                break;
            case "cluster":
                ClusterMediaLauncher.launch();
                break;
            default:
                System.out.println("Unknown command: " + command);
                System.out.println("Available commands: solo, cluster");
                break;
        }
    }
}
