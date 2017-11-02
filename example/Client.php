<?php

class CommanderClient {

    private $curlClient;

    private $socketPath;

    public function __construct(string $socketPath)
    {
        $this->curlClient = curl_init();
        $this->socketPath = $socketPath;

        curl_setopt($this->curlClient, CURLOPT_UNIX_SOCKET_PATH, $socketPath);
        curl_setopt($this->curlClient, CURLOPT_RETURNTRANSFER, true);
    }

    public function __destruct()
    {
        curl_close($this->curlClient);
    }

    private function generateRequestUri(string $requestPath)
    {
        /* Please note that Curl doesn't use http+unix:// or any other mechanism for
         *  specifying Unix Sockets; once the CURLOPT_UNIX_SOCKET_PATH option is set,
         *  Curl will simply ignore the domain of the request. Hence why this works,
         *  despite looking as though it should attempt to connect to a host found at
         *  the domain "unixsocket". See L14 where this is set.
         *
         *  @see Client.php:L14
         *  @see https://github.com/curl/curl/issues/1338
         */
        return sprintf("http://unixsocket%s", $requestPath);
    }


    public function dispatchCommand(string $command, array $parameters): array 
    {
        $payload = [
            "command"    => $command,
            "parameters" => $parameters
        ];

        curl_setopt($this->curlClient, CURLOPT_URL, "/dispatch");
        curl_setopt($this->curlClient, CURLOPT_POSTFIELDS, json_encode($payload));

        $result = curl_exec($this->curlClient);
        return json_decode($result, true);
    }


    public function listAvailableCommands(): array
    {
        curl_setopt($this->curlClient, CURLOPT_URL, $this->generateRequestUri("/listing"));

        $result = curl_exec($this->curlClient);
        return json_decode($result, true);
    }
}

$client = new CommanderClient('/tmp/commander.sock');
$availableCommands = $client->listAvailableCommands();

printf("Available Commander Commands (%s):\nName\t - Command\t - Description\n", $availableCommands['count']);
foreach ($availableCommands['commands'] as $command) {
    printf("%s\t - %s\t - %s\n", $command['name'], $command['command'], $command['description']);
}

