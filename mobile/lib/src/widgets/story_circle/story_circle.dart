import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';

class StoryCircle extends StatelessWidget {
  const StoryCircle({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 0, 0),
      child: CircleAvatar(
        radius: 35,
        backgroundColor: Colors.blue,
        child: CircleAvatar(
          radius: 32,
          backgroundColor: Theme.of(context).scaffoldBackgroundColor,
          child: CircleAvatar(
            backgroundImage: AssetImage(Constants.defaultAvatarPng),
            backgroundColor: Colors.white,
            radius: 30,
          ),
        ),
      ),
    );
  }
}
