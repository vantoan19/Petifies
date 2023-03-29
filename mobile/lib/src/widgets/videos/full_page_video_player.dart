// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';
import 'package:video_player/video_player.dart';

class FullPageVideoPlayer extends StatefulWidget {
  final String videoUrl;

  FullPageVideoPlayer({required this.videoUrl});

  @override
  _FullPageVideoPlayerState createState() => _FullPageVideoPlayerState();
}

class _FullPageVideoPlayerState extends State<FullPageVideoPlayer> {
  late VideoPlayerController _controller;
  bool _isPlaying = false;
  bool _isMuted = false;
  Duration _duration = Duration();
  Duration _position = Duration();
  double _sliderValue = 0.0;

  @override
  void initState() {
    super.initState();
    _controller = VideoPlayerController.network(widget.videoUrl)
      ..addListener(() {
        setState(() {
          _duration = _controller.value.duration;
          _position = _controller.value.position;
          _sliderValue = _position.inSeconds.toDouble();
        });
      })
      ..initialize().then((_) {
        setState(() {});
      });
  }

  void _togglePlaying() {
    setState(() {
      _isPlaying = !_isPlaying;
      if (_isPlaying) {
        _controller.play();
      } else {
        _controller.pause();
      }
    });
  }

  void _toggleMuted() {
    setState(() {
      _isMuted = !_isMuted;
      _controller.setVolume(_isMuted ? 0 : 1);
    });
  }

  void _changeSliderValue(value) {
    setState(() {
      _sliderValue = value;
      _controller.seekTo(Duration(seconds: _sliderValue.toInt()));
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.max,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Expanded(
          child: FittedBox(
            fit: BoxFit.fitWidth,
            child: Container(
              width: MediaQuery.of(context).size.width,
              child: AspectRatio(
                aspectRatio: _controller.value.aspectRatio,
                child: VideoPlayer(_controller),
              ),
            ),
          ),
        ),
        VideoControllerBar(
          isPlaying: _isPlaying,
          isMuted: _isMuted,
          pauseVideoCallback: _togglePlaying,
          muteVideoCallback: _toggleMuted,
          onSliderValueChanged: _changeSliderValue,
          sliderValue: _sliderValue,
          currentTime: _position,
          totalTime: _duration,
        ),
      ],
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
}

class VideoControllerBar extends StatelessWidget {
  final bool isPlaying;
  final bool isMuted;
  final VoidCallback pauseVideoCallback;
  final VoidCallback muteVideoCallback;
  final Function(double) onSliderValueChanged;
  final double sliderValue;
  final Duration currentTime;
  final Duration totalTime;

  const VideoControllerBar({
    Key? key,
    required this.isPlaying,
    required this.isMuted,
    required this.pauseVideoCallback,
    required this.muteVideoCallback,
    required this.onSliderValueChanged,
    required this.sliderValue,
    required this.currentTime,
    required this.totalTime,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Row(
        mainAxisSize: MainAxisSize.max,
        children: [
          IconButton(
            icon: Icon(
              isPlaying ? Icons.pause : Icons.play_arrow,
              color: Themes.whiteColor,
              size: 22,
            ),
            onPressed: pauseVideoCallback,
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8),
              child: SliderTheme(
                data: SliderTheme.of(context).copyWith(
                  overlayShape: SliderComponentShape.noThumb,
                  thumbColor: Themes.blueColor,
                  thumbShape: RoundSliderThumbShape(enabledThumbRadius: 1.5),
                ),
                child: Slider(
                  min: 0.0,
                  max: totalTime.inSeconds.toDouble(),
                  value: sliderValue.toDouble(),
                  onChanged: onSliderValueChanged,
                  activeColor: Themes.blueColor,
                  inactiveColor: Themes.lightGreyColor,
                ),
              ),
            ),
          ),
          Container(
            width: 50,
            child: Center(
              child: Text(
                '${currentTime.inMinutes}:${currentTime.inSeconds.remainder(60).toString().padLeft(2, '0')}',
                style: TextStyle(
                  color: Colors.white,
                ),
              ),
            ),
          ),
          IconButton(
            icon: Icon(
              isMuted ? Icons.volume_off : Icons.volume_up,
              color: Themes.whiteColor,
              size: 22,
            ),
            onPressed: muteVideoCallback,
          ),
        ],
      ),
    );
  }
}
