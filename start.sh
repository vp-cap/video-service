echo $1
if [ "$1" = "bo" ]
then
    echo ${GIT_USER}
    echo ${GIT_PASS}
    docker build --build-arg GIT_USER=${GIT_USER} --build-arg GIT_PASS=${GIT_PASS} -t video-service .
elif [ "$1" = "br" ]
then
    docker build --build-arg GIT_USER=${GIT_USER} --build-arg GIT_PASS=${GIT_PASS} -t video-service .
    docker stop video-service && docker rm video-service
    docker run --network=common --name video-service -p 50053:50053 video-service
else
    docker stop video-service && docker rm video-service
    docker run --network=common --name video-service -p 50053:50053 video-service
fi
