import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';

class TextDivider extends StatelessWidget {
  final String text;
  final double thickness;
  const TextDivider({super.key, required this.text, this.thickness = 1});

  @override
  Widget build(BuildContext context) {
    const textStyle = const TextStyle(
        color: Themes.greyColor, fontSize: 16, fontWeight: FontWeight.w300);

    return Row(
      children: [
        Expanded(
            child: Padding(
          padding: const EdgeInsets.symmetric(
              horizontal: Constants.horizontalScreenPadding, vertical: 8),
          child: Divider(
            thickness: thickness,
          ),
        )),
        Text(
          text,
          style: textStyle,
        ),
        Expanded(
            child: Padding(
          padding: const EdgeInsets.symmetric(
              horizontal: Constants.horizontalScreenPadding, vertical: 8),
          child: Divider(
            thickness: thickness,
          ),
        ))
      ],
    );
  }
}
