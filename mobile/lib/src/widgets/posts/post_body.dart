import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/videos/video.dart';
import 'package:video_player/video_player.dart';

class PostBody extends StatelessWidget {
  final String? textContent;
  final List<NetworkImageModel>? images;
  final List<NetworkVideoModel>? videos;

  const PostBody(
      {super.key,
      this.textContent = null,
      this.images = null,
      this.videos = null});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 18, 0, 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Text Content
          if (textContent != null)
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 0, 12, 18),
              child: Text(
                textContent!,
                style: TextStyle(
                  fontSize: 16,
                ),
              ),
            ),
          // Image & Video content
          if (images != null || videos != null)
            CarouselSlider(
              options: CarouselOptions(
                aspectRatio: 6 / 4,
                height: MediaQuery.of(context).size.width,
                enableInfiniteScroll: false,
                enlargeCenterPage: true,
                disableCenter: true,
                viewportFraction: 1.0,
              ),
              items: [
                // Images
                if (images != null)
                  ...images!
                      .map(
                        (image) => Builder(
                          builder: (BuildContext context) {
                            return Container(
                              width: MediaQuery.of(context).size.width,
                              decoration: BoxDecoration(
                                color: Themes.blackColor,
                              ),
                              child: AspectRatio(
                                aspectRatio: 4 / 6,
                                child: Image.network(
                                  image.uri,
                                  loadingBuilder:
                                      (context, child, loadingProgress) =>
                                          (loadingProgress == null)
                                              ? child
                                              : CircularProgressIndicator(),
                                  fit: BoxFit.fitWidth,
                                ),
                              ),
                            );
                          },
                        ),
                      )
                      .toList(),
                // Videos
                if (videos != null)
                  ...videos!.map(
                    (video) => Builder(
                      builder: (BuildContext context) {
                        return Container(
                          width: MediaQuery.of(context).size.width,
                          margin: EdgeInsets.symmetric(horizontal: 5.0),
                          child: VideoWidget(
                            videoPlayerController:
                                VideoPlayerController.network(
                              video.uri,
                            ),
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
      ),
    );
  }
}
