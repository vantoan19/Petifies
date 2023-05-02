// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';

class CustomTextFormField extends StatelessWidget {
  final String label;
  final Widget icon;
  final Widget? suffixIcon;
  final bool isObscureText;
  final Function(String) onChange;
  final bool? readOnly;
  final int? maxLines;
  final TextEditingController? controller;
  final String? Function(String?)? validator;

  const CustomTextFormField(
      {Key? key,
      required this.label,
      required this.icon,
      this.suffixIcon,
      this.readOnly,
      this.isObscureText = false,
      required this.onChange,
      this.controller = null,
      this.maxLines = 1,
      this.validator})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 8, horizontal: Constants.horizontalScreenPadding),
      child: TextFormField(
        obscureText: isObscureText,
        onChanged: onChange,
        validator: validator,
        controller: controller,
        readOnly: readOnly ?? false,
        maxLines: maxLines,
        decoration: InputDecoration(
          label: Text(label),
          prefixIcon: icon,
          suffixIcon: suffixIcon,
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
