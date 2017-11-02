<?php

class CommanderClient {

    /** @param resource */
    private $curlClient;

    /** @param string */
    private $socketPath;

    /** @param string|null */
    private $curlError = null;

    /**
     * Constructor: Initialises the Curl Resource, making it usable for subsequent
     *  API requests.
     *
     * @param string
     */
    public function __construct(string $socketPath)
    {
        $this->curlClient = curl_init();
        $this->socketPath = $socketPath;

        curl_setopt($this->curlClient, CURLOPT_UNIX_SOCKET_PATH, $socketPath);
        curl_setopt($this->curlClient, CURLOPT_RETURNTRANSFER, true);
    }

    /**
     * Deconstructor: Ensure the Curl Resource is correctly closed.
     */
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


    /**
     * Dispatches a command - via Curl - to Commander's Unix Socket.
     *
     * @param  string Command identifier to invoke on the Commander daemon.
     * @param  array  Parameters to dispatch to the Commander daemon.
     * @return array  JSON decoded response from Commander.
     */
    public function dispatchCommand(string $command, array $parameters): array 
    {
        $payload = [
            "command"    => $command,
            "parameters" => $parameters
        ];

        curl_setopt($this->curlClient, CURLOPT_URL, "/dispatch");
        curl_setopt($this->curlClient, CURLOPT_POSTFIELDS, json_encode($payload));

        $result = curl_exec($this->curlClient);
        if ($result === FALSE) {
            $this->curlError = curl_error($this->curlClient);
            return [ 'success' => false ];
        }

        return json_decode($result, true);
    }


    /**
     * Retrieves a list of all available commands from Commander's API.
     *
     * @return array  List of available commands.
     */
    public function listAvailableCommands(): array
    {
        curl_setopt($this->curlClient, CURLOPT_URL, $this->generateRequestUri("/listing"));

        $result = curl_exec($this->curlClient);
        if ($result === FALSE) {
            $this->curlError = curl_error($this->curlClient);
            return [ 'success' => false ];
        }
        
        return json_decode($result, true);
    }


    /**
     * Returns a human readable string from Curl in the event of an error.
     *
     * @return bool|string 
     */
    public function getCurlError()
    {
        return is_null($this->curlError) ? false : $this->curlError;
    }
}

$client = new CommanderClient('/tmp/commander.sock');
$availableCommands = $client->listAvailableCommands();

if (isset($availableCommands['status']) && $availableCommands['status'] === FALSE) {
    return printf("Error! (%s)\n", $client->getCurlError());
}

printf("Available Commander Commands (%s):\nName\t - Command\t - Description\n", $availableCommands['count']);
foreach ($availableCommands['commands'] as $command) {
    printf("%s\t - %s\t - %s\n", $command['name'], $command['command'], $command['description']);
}

