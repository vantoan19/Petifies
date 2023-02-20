import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';

class AuthTextField extends StatelessWidget {
  final String label;
  final Widget icon;
  final TextEditingController controller;
  final bool isObscureText;

  const AuthTextField(
      {super.key,
      required this.label,
      required this.icon,
      required this.controller,
      this.isObscureText = false});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 8, horizontal: Constants.horizontalScreenPadding),
      child: TextField(
        obscureText: isObscureText,
        controller: controller,
        decoration: InputDecoration(
          label: Text(label),
          prefixIcon: icon,
          enabledBorder: UnderlineInputBorder(
            borderSide: BorderSide(
                color: Theme.of(context).colorScheme.secondary, width: 1),
          ),
          errorBorder: const UnderlineInputBorder(
            borderSide: BorderSide(color: Themes.yellowColor, width: 1),
          ),
          focusedBorder: UnderlineInputBorder(
            borderSide: BorderSide(
                color: Theme.of(context).colorScheme.primary, width: 2.5),
          ),
          focusedErrorBorder: const UnderlineInputBorder(
            borderSide: BorderSide(color: Themes.redColor, width: 2.5),
          ),
        ),
        style: Theme.of(context).textTheme.bodyMedium,
      ),
    );
  }
}
