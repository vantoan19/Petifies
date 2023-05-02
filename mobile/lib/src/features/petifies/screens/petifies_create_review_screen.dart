// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter_datetime_picker/flutter_datetime_picker.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/review_controllers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';

class CreateReviewScreen extends ConsumerStatefulWidget {
  final String petifiesId;
  const CreateReviewScreen({
    Key? key,
    required this.petifiesId,
  }) : super(key: key);

  @override
  ConsumerState<CreateReviewScreen> createState() => _CreateReviewScreenState();
}

class _CreateReviewScreenState extends ConsumerState<CreateReviewScreen> {
  TextEditingController _reviewTextController = TextEditingController();

  @override
  void initState() {
    super.initState();
  }

  Future<void> _submitCreation() async {
    try {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(
          "Creating review",
        ),
      ));
      final review =
          await ref.read(createReviewControllerProvider.notifier).createReview(
                petifiesId: widget.petifiesId,
                review: _reviewTextController.text,
              );
      ref
          .read(listReviewsByPetifidesIdControllerProvider(
                  petifiesId: widget.petifiesId)
              .notifier)
          .addReview(review);
      ScaffoldMessenger.of(context).hideCurrentSnackBar();
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(
        "Review created successfully",
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
                "New Review",
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
            TextFormField(
              autofocus: true,
              controller: _reviewTextController,
              decoration: InputDecoration(
                label: Text(
                  "Review",
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                hintText: 'Write your review...',
                contentPadding: EdgeInsets.symmetric(vertical: 10),
              ),
              style: TextStyle(
                fontSize: 15,
              ),
              maxLines: null,
              keyboardType: TextInputType.multiline,
            ),
            Padding(
              padding: const EdgeInsets.only(top: 32.0),
              child: ElevatedButton(
                onPressed: _submitCreation,
                style: ElevatedButton.styleFrom(
                  minimumSize: Size.fromHeight(44),
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(50)),
                ),
                child: Text(
                  "Create Review",
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
