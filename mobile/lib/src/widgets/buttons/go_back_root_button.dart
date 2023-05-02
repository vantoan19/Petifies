import 'package:flutter/material.dart';
import 'package:mobile/src/utils/navigation.dart';

class GoBackRootButton extends StatelessWidget {
  const GoBackRootButton({super.key});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () => {NavigatorUtil.goBackRootNavigater(context)},
      icon: Icon(Icons.arrow_back),
    );
  }
}
