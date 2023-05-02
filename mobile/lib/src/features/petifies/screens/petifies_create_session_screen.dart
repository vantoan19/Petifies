// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';
import 'package:flutter_datetime_picker/flutter_datetime_picker.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_session_controllers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';
import 'package:mobile/src/widgets/buttons/pick_datetime_button.dart';

class PetifiesCreateSessionScreen extends ConsumerStatefulWidget {
  final String petifiesId;
  const PetifiesCreateSessionScreen({
    Key? key,
    required this.petifiesId,
  }) : super(key: key);

  @override
  ConsumerState<PetifiesCreateSessionScreen> createState() =>
      _PetifiesCreateSessionScreenState();
}

class _PetifiesCreateSessionScreenState
    extends ConsumerState<PetifiesCreateSessionScreen> {
  DateTime? _fromTime;
  DateTime? _toTime;
  TextEditingController _fromTimeTextController =
      TextEditingController(text: "DD/MM/YYYY HH:MM");
  TextEditingController _toTimeTextController =
      TextEditingController(text: "DD/MM/YYYY HH:MM");

  @override
  void initState() {
    super.initState();
  }

  String digits(int value, int length) {
    return '$value'.padLeft(length, "0");
  }

  String convertDateTime(DateTime time) {
    return '${digits(time.day, 2)}/${digits(time.month, 2)}/${digits(time.year, 4)} ${digits(time.hour, 2)}:${digits(time.minute, 2)}';
  }

  Future<void> _submitCreation() async {
    try {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(
          "Creating session",
        ),
      ));

      final session = await ref
          .read(createPetifiesSessionControllerProvider.notifier)
          .createPetifiesSession(
            petifiesId: widget.petifiesId,
            fromTime: _fromTime!,
            toTime: _toTime!,
          );
      ref
          .read(listSessionsByPetifiesIdControllerProvider(
                  petifiesId: widget.petifiesId)
              .notifier)
          .addSession(session);
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(
        "Session created successfully",
      )));

      NavigatorUtil.goBack(context);
    } on GrpcError catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(
        e.message ?? e.toString(),
      )));
    } catch (e) {
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(
        e.toString(),
      )));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: OnlyGoBackAppbar(),
      body: Padding(
        padding: const EdgeInsets.symmetric(
            horizontal: Constants.petifiesExpoloreHorizontalPadding),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 28, 0, 20),
              child: Text(
                "New Session",
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
            _SessionFormField(
              fieldLabel: "From time: ",
              onConfirm: (date) {
                setState(() {
                  _fromTime = date;
                  _fromTimeTextController.text = convertDateTime(date);
                });
              },
              controller: _fromTimeTextController,
            ),
            _SessionFormField(
              fieldLabel: "To time: ",
              onConfirm: (date) {
                setState(() {
                  _toTime = date;
                  _toTimeTextController.text = convertDateTime(date);
                });
              },
              controller: _toTimeTextController,
            ),
            Padding(
              padding: const EdgeInsets.only(top: 32.0),
              child: ElevatedButton(
                onPressed: (_toTime == null || _fromTime == null)
                    ? null
                    : () => _submitCreation(),
                style: ElevatedButton.styleFrom(
                  minimumSize: Size.fromHeight(44),
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(50)),
                ),
                child: Text(
                  "Create Session",
                  style: TextStyle(
                    color: Themes.whiteColor,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _SessionFormField extends StatelessWidget {
  final String fieldLabel;
  final Function(DateTime) onConfirm;
  final TextEditingController controller;
  const _SessionFormField({
    Key? key,
    required this.fieldLabel,
    required this.onConfirm,
    required this.controller,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.max,
      children: [
        Expanded(
          child: TextFormField(
            controller: controller,
            style: TextStyle(
              color: controller.text == "DD/MM/YYYY HH:MM"
                  ? Theme.of(context).colorScheme.secondary.withOpacity(0.3)
                  : Theme.of(context).colorScheme.secondary,
              fontSize: 15,
              height: 1.25,
            ),
            decoration: InputDecoration(
              label: Text(
                fieldLabel,
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w800,
                ),
              ),
            ),
            readOnly: true,
            onTap: () {
              DatePicker.showDateTimePicker(
                context,
                showTitleActions: true,
                minTime: DateTime.now(),
                onConfirm: onConfirm,
                currentTime: DateTime.now(),
              );
            },
          ),
        ),
      ],
    );
  }
}
