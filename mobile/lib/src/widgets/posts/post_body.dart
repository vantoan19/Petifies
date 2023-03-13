import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/widgets/videos/video.dart';
import 'package:video_player/video_player.dart';

class PostBody extends StatelessWidget {
  final String? textContent;
  final List<ImageModel>? images;
  final List<VideoModel>? videos;

  const PostBody(
      {super.key,
      this.textContent = null,
      this.images = null,
      this.videos = null});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        if (textContent != null)
          Text(
            textContent!,
            style: TextStyle(
              fontSize: 11,
            ),
          ),
        if (images != null || videos != null)
          CarouselSlider(
            options: CarouselOptions(
              aspectRatio: 4 / 6,
              height: 400,
              enableInfiniteScroll: false,
            ),
            items: [
              if (images != null)
                ...images!
                    .map(
                      (image) => Builder(
                        builder: (BuildContext context) {
                          return Container(
                            width: MediaQuery.of(context).size.width,
                            margin: EdgeInsets.symmetric(horizontal: 5.0),
                            child: AspectRatio(
                              aspectRatio: 4 / 6,
                              child: Image.network(image.uri),
                            ),
                          );
                        },
                      ),
                    )
                    .toList(),
              if (videos != null)
                ...videos!.map(
                  (video) => Builder(
                    builder: (BuildContext context) {
                      return Container(
                        width: MediaQuery.of(context).size.width,
                        margin: EdgeInsets.symmetric(horizontal: 5.0),
                        child: VideoWidget(
                          videoPlayerController:
                              VideoPlayerController.network(video.uri),
                          isAutoplaying: true,
                          isLooping: true,
                        ),
                      );
                    },
                  ),
                )
            ],
          )
      ],
    );
  }
}
