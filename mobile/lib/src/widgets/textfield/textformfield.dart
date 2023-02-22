
import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';

class CustomTextFormField extends StatelessWidget {
  final String label;
  final Widget icon;
  final bool isObscureText;
  final Function(String) onChange;
  final String? Function(String?)? validator;

  const CustomTextFormField(
      {super.key,
      required this.label,
      required this.icon,
      this.isObscureText = false,
      required this.onChange,
      this.validator});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 8, horizontal: Constants.horizontalScreenPadding),
      child: TextFormField(
        obscureText: isObscureText,
        onChanged: onChange,
        validator: validator,
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
