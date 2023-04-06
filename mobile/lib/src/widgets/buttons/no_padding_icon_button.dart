import 'package:flutter/material.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';

class NoPaddingIconButton extends IconButton {
  const NoPaddingIconButton({
    super.key,
    required super.onPressed,
    required super.icon,
    super.iconSize,
    super.padding = EdgeInsets.zero,
    super.constraints = const BoxConstraints(
      minHeight: 15,
      minWidth: 15,
      maxHeight: 40,
      maxWidth: 40,
    ),
    super.color,
  });
}
