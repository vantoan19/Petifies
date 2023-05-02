// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';

class IconWithLabelTab extends StatelessWidget {
  final Icon icon;
  final String label;
  const IconWithLabelTab({
    Key? key,
    required this.icon,
    required this.label,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Tab(
      icon: icon,
      text: label,
    );
  }
}
