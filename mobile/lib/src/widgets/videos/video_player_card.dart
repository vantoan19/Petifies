// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/media/screens/media_full_page_screen.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:video_player/video_player.dart';

class VideoPlayerCard extends ConsumerStatefulWidget {
  final bool isRoundedTopLeft;
  final bool isRoundedTopRight;
  final bool isRoundedBottomLeft;
  final bool isRoundedBottomRight;
  final String videoUrl;
  final bool autoPlay;
  final double? height; // If height is specified, maxHeight will be height
  final double width;
  final double maxHeight;
  final EdgeInsetsGeometry? padding;
  final VideoPlayerController? controller;
  final VoidCallback playNextVideoCallback;

  final bool isClickable;

  const VideoPlayerCard({
    Key? key,
    required this.isRoundedTopLeft,
    required this.isRoundedTopRight,
    required this.isRoundedBottomLeft,
    required this.isRoundedBottomRight,
    required this.videoUrl,
    required this.autoPlay,
    this.height,
    required this.width,
    required this.maxHeight,
    this.padding,
    this.controller = null,
    required this.playNextVideoCallback,
    required this.isClickable,
  }) : super(key: key);

  @override
  ConsumerState<VideoPlayerCard> createState() => _VideoPlayerCardState();
}

class _VideoPlayerCardState extends ConsumerState<VideoPlayerCard> {
  late VideoPlayerController _controller;
  bool _isMuted = false;
  Duration _duration = Duration();
  Duration _position = Duration();

  @override
  void initState() {
    super.initState();
    if (widget.controller != null) {
      _controller = widget.controller!;
    } else {
      _controller = VideoPlayerController.network(widget.videoUrl);
    }
    _controller
      ..addListener(() {
        setState(() {
          _duration = _controller.value.duration;
          _position = _controller.value.position;
          if (_controller.value.isInitialized &&
              (_controller.value.duration == _controller.value.position)) {
            //checking the duration and position every time
            widget.playNextVideoCallback();
          }
        });
      })
      ..initialize().then((_) {
        setState(() {
          if (widget.autoPlay) {
            _controller.play();
          }
        });
      });
  }

  void _toggleMuted() {
    setState(() {
      _isMuted = !_isMuted;
      _controller.setVolume(_isMuted ? 0 : 1);
    });
  }

  @override
  void didUpdateWidget(covariant VideoPlayerCard oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.autoPlay != widget.autoPlay && widget.autoPlay) {
      _controller.play();
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final remainingTime = _duration - _position;

    double maxHeight;
    if (widget.height != null) {
      maxHeight = widget.height!;
    } else {
      maxHeight = _controller.value.isInitialized
          ? min(
              widget.maxHeight,
              widget.width *
                  _controller.value.size.height /
                  _controller.value.size.width,
            )
          : widget.maxHeight;
    }

    Widget videoCard = Container(
      width: widget.width,
      height: widget.height,
      padding: widget.padding,
      constraints: BoxConstraints(
        maxHeight: maxHeight,
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.only(
          topLeft: widget.isRoundedTopLeft ? Radius.circular(12) : Radius.zero,
          topRight:
              widget.isRoundedTopRight ? Radius.circular(12) : Radius.zero,
          bottomLeft:
              widget.isRoundedBottomLeft ? Radius.circular(12) : Radius.zero,
          bottomRight:
              widget.isRoundedBottomRight ? Radius.circular(12) : Radius.zero,
        ),
        child: Stack(
          alignment: Alignment.bottomCenter,
          fit: StackFit.expand,
          children: [
            FittedBox(
              fit: BoxFit.cover,
              child: SizedBox(
                width: _controller.value.aspectRatio,
                height: 1,
                child: VideoPlayer(_controller),
              ),
            ),
            Align(
              alignment: Alignment.bottomCenter,
              child: Container(
                height: 50.0,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Container(
                          height: 30,
                          width: 60,
                          margin: EdgeInsets.fromLTRB(12, 0, 0, 0),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.all(Radius.circular(18)),
                            color: Colors.black54,
                          ),
                          child: Center(
                            child: Text(
                              "${remainingTime.inMinutes}:${remainingTime.inSeconds.remainder(60).toString().padLeft(2, '0')}",
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 14,
                                fontWeight: FontWeight.w700,
                              ),
                            ),
                          ),
                        ),
                        Container(
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: Colors.black54,
                          ),
                          margin: EdgeInsets.fromLTRB(0, 0, 12, 0),
                          child: IconButton(
                            padding: EdgeInsets.all(0),
                            constraints: BoxConstraints(
                              minHeight: 30,
                              minWidth: 30,
                            ),
                            icon: Icon(
                              _isMuted ? Icons.volume_off : Icons.volume_up,
                              color: Colors.white,
                              size: 22,
                            ),
                            onPressed: _toggleMuted,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );

    if (!widget.isClickable) {
      return videoCard;
    }

    return GestureDetector(
      onTap: () {
        showModalBottomSheet(
            context: context,
            isScrollControlled: true,
            useSafeArea: true,
            barrierColor: Theme.of(context).scaffoldBackgroundColor,
            builder: (context) {
              return NavigatorUtil.showMediaFullPageBottomSheet(
                ref: ref,
                mediaUrl: widget.videoUrl,
                isMediaImage: false,
              );
            });
      },
      child: videoCard,
    );
  }
}
