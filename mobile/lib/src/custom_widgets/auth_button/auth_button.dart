import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';

class AuthButton extends StatelessWidget {
  final String label;
  final Function action;
  final Color color;

  const AuthButton(
      {super.key,
      required this.label,
      required this.action,
      this.color = Themes.blueColor});

  const AuthButton.withColor(
      {super.key,
      required this.label,
      required this.action,
      required this.color});

  Text get _label {
    return Text(label);
  }

  Widget get _icon {
    return const SizedBox();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          horizontal: Constants.horizontalScreenPadding, vertical: 8),
      child: ElevatedButton.icon(
        onPressed: () => {action()},
        icon: _icon,
        label: _label,
        style: ElevatedButton.styleFrom(
            backgroundColor: color,
            minimumSize: const Size(double.infinity, 50.0),
            shape:
                RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
            textStyle: Theme.of(context).textTheme.titleMedium),
      ),
    );
  }
}

class ThirdpartyAuthButton extends AuthButton {
  final Widget icon;

  const ThirdpartyAuthButton(
      {super.key,
      required super.label,
      required super.action,
      required this.icon});

  const ThirdpartyAuthButton.withColor(
      {super.key,
      required super.label,
      required super.action,
      required super.color,
      required this.icon});

  @override
  Widget get _icon {
    return icon;
  }
}
