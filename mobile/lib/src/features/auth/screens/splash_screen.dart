import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:video_player/video_player.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  late VideoPlayerController _controller;

  @override
  void initState() {
    super.initState();

    _controller = VideoPlayerController.asset(Constants.splashIconMp4Path)
      ..initialize().then((_) {
        setState(() {});
      })
      ..setVolume(0.0);

    _play();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _play() async {
    _controller.play();

    await Future.delayed(const Duration(milliseconds: 1500));

    Navigator.of(context).pushNamed("/introduction");
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Themes.splashBlueColor,
      body: Center(
        child: _controller.value.isInitialized
            ? AspectRatio(
                aspectRatio: _controller.value.aspectRatio,
                child: VideoPlayer(
                  _controller,
                ),
              )
            : Container(),
      ),
    );
  }
}
