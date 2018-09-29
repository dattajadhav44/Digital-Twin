Using dev mode
==============

The stpes are bit lengthy inorder to execute this exrcize. Please do carefully fallow with below steps. and then invoke the chaincode dcar.

mainly we have 5 functions in our dcar chaince code. Let me brief abouth them.

1.initLedger:- This function initialize the ledger and 5 car entries for us in database(ledger).

2.createCar:- This functions create the car. basically adding new entry into ledger.

3.changeCarOwner:- This function changes the ownership of the car based in VIN (Vehicle identification number)

4.queryAllCars:- The function is fetch all the cars available in range of 0 to 10.

5.queryCar:- This function is to fetch the specific based on VIN number. 

6.changeCarMilleageAndColour:- This function is to modify the milleage and colour of the car based on VIN number.

The queries are given below with all the parameters. Please see the below.
----------------------------------------------------------------------------------------------

Install Fabric Samples
----------------------

If you haven't already done so, please install the doc [samples](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html).

Navigate to the ``chaincode-docker-devmode`` directory of the ``fabric-samples``
clone:

  cd chaincode-docker-devmode

Download docker images
^^^^^^^^^^^^^^^^^^^^^^

We need four docker images in order for "dev mode" to run against the supplied
docker compose script.  If you installed the ``fabric-samples`` repo clone and
followed the instructions to [download-platform-specific-binaries](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#download-platform-specific-binaries), then
you should have the necessary Docker images installed locally.

.. note:: If you choose to manually pull the images then you must retag them as
          ``latest``.

Issue a ``docker images`` command to reveal your local Docker Registry.  You
should see something similar to following:

  docker images
  REPOSITORY                     TAG                                  IMAGE ID            CREATED             SIZE
  hyperledger/fabric-tools       latest                c584c20ac82b        9 days ago         1.42 GB
  hyperledger/fabric-tools       x86_64-1.1.0-preview  c584c20ac82b        9 days ago         1.42 GB
  hyperledger/fabric-orderer     latest                2fccc91736df        9 days ago         159 MB
  hyperledger/fabric-orderer     x86_64-1.1.0-preview  2fccc91736df        9 dyas ago         159 MB
  hyperledger/fabric-peer        latest                337f3d90b452        9 days ago         165 MB
  hyperledger/fabric-peer        x86_64-1.1.0-preview  337f3d90b452        9 days ago         165 MB
  hyperledger/fabric-ccenv       latest                82489d1c11e8        9 days ago         1.35 GB
  hyperledger/fabric-ccenv       x86_64-1.1.0-preview  82489d1c11e8        9 days ago         1.35 GB

.. note:: If you retrieved the images through the [download-platform-specific-binaries](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#download-platform-specific-binaries),
          then you will see additional images listed.  However, we are only concerned with
          these four.

Now open three terminals and navigate to your ``chaincode-docker-devmode``
directory in each.

Terminal 1 - Start the network
------------------------------

.. code:: bash

    docker-compose -f docker-compose-simple.yaml up

The above starts the network with the ``SingleSampleMSPSolo`` orderer profile and
launches the peer in "dev mode".  It also launches two additional containers -
one for the chaincode environment and a CLI to interact with the chaincode.  The
commands for create and join channel are embedded in the CLI container, so we
can jump immediately to the chaincode calls.

Terminal 2 - Build & start the chaincode
----------------------------------------

.. code:: bash

  docker exec -it chaincode bash

You should see the following:

.. code:: bash

  root@d2629980e76b:/opt/gopath/src/chaincode#

Now, compile your chaincode:

.. code:: bash

  cd dcar
  go build -o dcar

Now run the chaincode:

.. code:: bash


  CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./dcar

The chaincode is started with peer and chaincode logs indicating successful registration with the peer.
Note that at this stage the chaincode is not associated with any channel. This is done in subsequent steps
using the ``instantiate`` command.

Terminal 3 - Use the chaincode
------------------------------

Even though you are in ``--peer-chaincodedev`` mode, you still have to install the
chaincode so the life-cycle system chaincode can go through its checks normally.
This requirement may be removed in future when in ``--peer-chaincodedev`` mode.

We'll leverage the CLI container to drive these calls.


  docker exec -it cli bash

Install the chaincode here
  peer chaincode install -p chaincodedev/chaincode/dcar -n mycc -v 0

instantiate the chaincode here 

peer chaincode instantiate -n mycc -v 0 -c '{"Args":["init"," "]}' -C myc
  
peer chaincode invoke -n mycc -c '{"Args":["initLedger", " "]}' -C myc   - This will create 5 bran new entries for us. file records
 
peer chaincode invoke -n mycc -c '{"Args":["queryAllCars", " "]}' -C myc
 
peer chaincode invoke -n mycc -c '{"Args":["changeCarOwner", "3VW5DAAT6JM516495, Stephen"]}' -C myc
 
peer chaincode invoke -n mycc -c '{"Args":["changeCarMilleageAndColour", "WBS8M9C51J5K98915","green","30]}' -C myc
 
peer chaincode invoke -n mycc -c '{"Args":["createCar", "WDAPF4CC2JP603170", "Mathew", "grey", "Sprinter", "Mercedes-Benz", "25"]}' -C myc


.. code:: bash

  peer chaincode invoke -n mycc -c '{"Args":["queryAllCars", " "]}' -C myc
