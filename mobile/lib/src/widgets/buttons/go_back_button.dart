import 'package:flutter/material.dart';
import 'package:mobile/src/utils/navigation.dart';

class GoBackButton extends StatelessWidget {
  const GoBackButton({super.key});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () => {NavigatorUtil.goBack(context)},
      icon: Icon(Icons.arrow_back),
    );
  }
}
