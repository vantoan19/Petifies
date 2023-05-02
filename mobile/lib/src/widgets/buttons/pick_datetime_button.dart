// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';
import 'package:flutter_datetime_picker/flutter_datetime_picker.dart';

class PickDatetimeButton extends StatelessWidget {
  final String label;
  final Function(DateTime) onConfirm;
  const PickDatetimeButton({
    Key? key,
    required this.label,
    required this.onConfirm,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return OutlinedButton(
      onPressed: () {
        DatePicker.showDateTimePicker(
          context,
          showTitleActions: true,
          minTime: DateTime.now(),
          onConfirm: onConfirm,
          currentTime: DateTime.now(),
        );
      },
      child: Text(label),
    );
  }
}
