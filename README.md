# spbstu-smart-library
Smart analytics of data from SPBSTU library

# Startup process

Open cmd as Administrator and go to the root directory of distributive

Run the following commands from the root directory to install the whole package:

##### Note: First run can take a time

    cd ./bin
    install_all.bat

After that wait for opening Kibana GUI in the browser. Downloading of data will be processed in the background after the completion of the script. It may take a while. When data will be downloaded you can see them in http://.../app/discover#/ and find graphs here http://.../app/visualize#/

If all software already installed you can startup system without installing software via the following commands:

    cd ./bin
    setup_ek.bat


# Cleanup

To remove the package (with third-party software) open cmd as Administrator, go to the root directory and run the following commands:

    cd ./bin
    cleanup_all.bat

To remove dockers only run the following commands:

    cd ./bin
    cleanup_docker.bat
