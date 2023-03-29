// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/widgets/images/image_card.dart';
import 'package:mobile/src/widgets/videos/video_player_card.dart';
import 'package:video_player/video_player.dart';

class MediaView extends StatefulWidget {
  final List<String> imageUrls;
  final List<File>? imageFiles;
  final List<String> videoUrls;
  final List<VideoPlayerController>? videoControllers;

  final bool isClickable;

  const MediaView({
    Key? key,
    required this.imageUrls,
    this.imageFiles = null,
    required this.videoUrls,
    this.videoControllers = null,
    required this.isClickable,
  }) : super(key: key);

  @override
  State<MediaView> createState() => _MediaViewState();
}

class _MediaViewState extends State<MediaView> {
  int _currentPlayingVideoIndex = 0;

  void _playNextVideo() {
    setState(() {
      _currentPlayingVideoIndex = _currentPlayingVideoIndex + 1;
    });
  }

  @override
  Widget build(BuildContext context) {
    final imagesLength = (widget.imageFiles != null
        ? widget.imageFiles!.length
        : widget.imageUrls.length);
    final videosLength = (widget.videoControllers != null
        ? widget.videoControllers!.length
        : widget.videoUrls.length);
    final mediaLength = imagesLength + videosLength;

    Widget grid;

    switch (mediaLength) {
      // ==============================================
      // ================ Only 1 media ================
      // ==============================================
      case 1:
        final mediaWidth = MediaQuery.of(context).size.width;
        final maxMediaHeight = mediaWidth * 1.2;
        grid = (imagesLength > 0)
            ? ImageCard(
                isRoundedTopLeft: true,
                isRoundedTopRight: true,
                isRoundedBottomLeft: true,
                isRoundedBottomRight: true,
                width: mediaWidth,
                maxHeight: maxMediaHeight,
                padding: EdgeInsets.symmetric(
                    horizontal: Constants.horizontalScreenPadding),
                imageUrl:
                    widget.imageUrls.length > 0 ? widget.imageUrls[0] : "",
                imageFile:
                    (widget.imageFiles != null && widget.imageFiles!.length > 0)
                        ? widget.imageFiles![0]
                        : null,
                isClickable: widget.isClickable,
              )
            : VideoPlayerCard(
                isRoundedTopLeft: true,
                isRoundedTopRight: true,
                isRoundedBottomLeft: true,
                isRoundedBottomRight: true,
                width: mediaWidth,
                maxHeight: maxMediaHeight,
                padding: EdgeInsets.symmetric(
                    horizontal: Constants.horizontalScreenPadding),
                videoUrl:
                    widget.videoUrls.length > 0 ? widget.videoUrls[0] : "",
                controller: (widget.videoControllers != null &&
                        widget.videoControllers!.length > 0)
                    ? widget.videoControllers![0]
                    : null,
                autoPlay: _currentPlayingVideoIndex == 0,
                playNextVideoCallback: _playNextVideo,
                isClickable: widget.isClickable,
              );
        break;
      case 2:
        // ==============================================
        // ================== 2 medias ==================
        // ==============================================
        final mediaWidth = (MediaQuery.of(context).size.width -
                Constants.horizontalScreenPadding * 2) /
            2;
        final mediaHeight = ((MediaQuery.of(context).size.width -
                    Constants.horizontalScreenPadding * 2) /
                2) *
            1.4;

        List<Widget> children = [];
        for (int i = 0; i < 2; i++) {
          if (i >= videosLength) {
            children.add(ImageCard(
              isRoundedTopLeft: i == 0,
              isRoundedTopRight: i == 1,
              isRoundedBottomLeft: i == 0,
              isRoundedBottomRight: i == 1,
              width: mediaWidth,
              height: mediaHeight,
              maxHeight: mediaHeight,
              padding: EdgeInsets.symmetric(horizontal: 3),
              imageUrl: widget.imageUrls.length > 0
                  ? widget.imageUrls[i - videosLength]
                  : "",
              imageFile:
                  (widget.imageFiles != null && widget.imageFiles!.length > 0)
                      ? widget.imageFiles![i - videosLength]
                      : null,
              isClickable: widget.isClickable,
            ));
          } else {
            children.add(VideoPlayerCard(
              isRoundedTopLeft: i == 0,
              isRoundedTopRight: i == 1,
              isRoundedBottomLeft: i == 0,
              isRoundedBottomRight: i == 1,
              width: mediaWidth,
              height: mediaHeight,
              maxHeight: mediaHeight,
              padding: EdgeInsets.symmetric(horizontal: 3),
              videoUrl: widget.videoUrls.length > 0 ? widget.videoUrls[i] : "",
              controller: (widget.videoControllers != null &&
                      widget.videoControllers!.length > 0)
                  ? widget.videoControllers![i]
                  : null,
              autoPlay: _currentPlayingVideoIndex == (i),
              playNextVideoCallback: _playNextVideo,
              isClickable: widget.isClickable,
            ));
          }
        }
        grid = Container(
          padding: EdgeInsets.symmetric(
              horizontal: Constants.horizontalScreenPadding - 3),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: children,
          ),
        );
        break;
      case 3:
        // ==============================================
        // ================== 3 medias ==================
        // ==============================================
        final mediaWidth = (MediaQuery.of(context).size.width -
                Constants.horizontalScreenPadding * 2) /
            2;
        final mediaHeight = ((MediaQuery.of(context).size.width -
                    Constants.horizontalScreenPadding * 2) /
                2) *
            1.4;
        List<Widget> children = [];
        for (int i = 0; i < 3; i++) {
          if (i >= videosLength) {
            children.add(ImageCard(
              isRoundedTopLeft: i == 0,
              isRoundedTopRight: i == 1,
              isRoundedBottomLeft: i == 0,
              isRoundedBottomRight: i == 2,
              width: mediaWidth,
              height: i == 0 ? mediaHeight : mediaHeight / 2,
              maxHeight: i == 0 ? mediaHeight : mediaHeight / 2,
              padding: EdgeInsets.fromLTRB(
                3,
                i == 2 ? 3 : 0,
                3,
                i == 1 ? 3 : 0,
              ),
              imageUrl: widget.imageUrls.length > 0
                  ? widget.imageUrls[i - videosLength]
                  : "",
              imageFile:
                  (widget.imageFiles != null && widget.imageFiles!.length > 0)
                      ? widget.imageFiles![i - videosLength]
                      : null,
              isClickable: widget.isClickable,
            ));
          } else {
            children.add(VideoPlayerCard(
              isRoundedTopLeft: i == 0,
              isRoundedTopRight: i == 1,
              isRoundedBottomLeft: i == 0,
              isRoundedBottomRight: i == 2,
              width: mediaWidth,
              height: i == 0 ? mediaHeight : mediaHeight / 2,
              maxHeight: i == 0 ? mediaHeight : mediaHeight / 2,
              padding: EdgeInsets.fromLTRB(
                3,
                i == 2 ? 3 : 0,
                3,
                i == 1 ? 3 : 0,
              ),
              videoUrl: widget.videoUrls.length > 0 ? widget.videoUrls[i] : "",
              controller: (widget.videoControllers != null &&
                      widget.videoControllers!.length > 0)
                  ? widget.videoControllers![i]
                  : null,
              autoPlay: _currentPlayingVideoIndex == (i),
              playNextVideoCallback: _playNextVideo,
              isClickable: widget.isClickable,
            ));
          }
        }
        grid = Container(
          padding: EdgeInsets.symmetric(
              horizontal: Constants.horizontalScreenPadding),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              children[0],
              Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [...children.sublist(1)],
              )
            ],
          ),
        );
        break;
      case 4:
        // ==============================================
        // ================== 4 medias ==================
        // ==============================================
        final mediaWidth = (MediaQuery.of(context).size.width -
                Constants.horizontalScreenPadding * 2) /
            2;
        final mediaHeight = ((MediaQuery.of(context).size.width -
                    Constants.horizontalScreenPadding * 2) /
                2) *
            0.7;

        List<Widget> children = [];
        for (int i = 0; i < 4; i++) {
          if (i >= videosLength) {
            children.add(ImageCard(
              isRoundedTopLeft: i == 0,
              isRoundedTopRight: i == 1,
              isRoundedBottomLeft: i == 2,
              isRoundedBottomRight: i == 3,
              width: mediaWidth,
              height: mediaHeight,
              maxHeight: mediaHeight,
              padding: EdgeInsets.all(3),
              imageUrl: widget.imageUrls.length > 0
                  ? widget.imageUrls[i - videosLength]
                  : "",
              imageFile:
                  (widget.imageFiles != null && widget.imageFiles!.length > 0)
                      ? widget.imageFiles![i - videosLength]
                      : null,
              isClickable: widget.isClickable,
            ));
          } else {
            children.add(VideoPlayerCard(
              isRoundedTopLeft: i == 0,
              isRoundedTopRight: i == 1,
              isRoundedBottomLeft: i == 2,
              isRoundedBottomRight: i == 3,
              width: mediaWidth,
              height: mediaHeight,
              maxHeight: mediaHeight,
              padding: EdgeInsets.all(3),
              videoUrl: widget.videoUrls.length > 0 ? widget.videoUrls[i] : "",
              controller: (widget.videoControllers != null &&
                      widget.videoControllers!.length > 0)
                  ? widget.videoControllers![i]
                  : null,
              autoPlay: _currentPlayingVideoIndex == (i),
              playNextVideoCallback: _playNextVideo,
              isClickable: widget.isClickable,
            ));
          }
        }
        grid = Container(
          padding: EdgeInsets.symmetric(
              horizontal: Constants.horizontalScreenPadding),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [children[0], children[2]],
              ),
              Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [children[1], children[3]],
              )
            ],
          ),
        );
        break;
      default:
        grid = SizedBox.shrink();
        break;
    }
    return grid;
  }
}
