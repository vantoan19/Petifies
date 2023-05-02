// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:mobile/src/features/petifies/screens/list_proposals_for_session_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_proposal_screen.dart';
import 'package:mobile/src/models/petifies_session.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/status/status.dart';

class PetifiesSession extends ConsumerWidget {
  final bool showListProposalButton;
  const PetifiesSession({
    required this.showListProposalButton,
  });

  Widget statusWidget(PetifiesSessionStatus status) {
    switch (status) {
      case PetifiesSessionStatus.WAITING_FOR_PROPOSAL:
        return Status(
          color: Themes.blueColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "WAITING FOR PROPOSALS",
          textSize: 8,
          textColor: Themes.whiteColor,
        );
      case PetifiesSessionStatus.PROPOSAL_ACCEPTED:
        return Status(
          color: Themes.greenColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "PROPOSAL ACCEPTED",
          textSize: 8,
          textColor: Themes.whiteColor,
        );
      case PetifiesSessionStatus.ON_GOING:
        return Status(
          color: Themes.yellowColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "ON GOING",
          textSize: 8,
          textColor: Themes.whiteColor,
        );
      default:
        return Status(
          color: Themes.greyColor,
          margin: EdgeInsets.symmetric(vertical: 8),
          label: "UNKNOWN",
          textSize: 8,
          textColor: Themes.blackColor,
        );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final id = ref.watch(petifiesSessionProvider.select((value) => value.id));
    final status =
        ref.watch(petifiesSessionProvider.select((value) => value.status));
    final fromTime =
        ref.watch(petifiesSessionProvider.select((value) => value.fromTime));
    final toTime =
        ref.watch(petifiesSessionProvider.select((value) => value.toTime));
    final DateFormat formatter = DateFormat("HH:mm dd-MM-yyyy");

    final user = ref.watch(myUserProvider);
    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    return Card(
      color: Theme.of(context).colorScheme.tertiary,
      child: Container(
        width: 250,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(16.0, 16, 16, 0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ListTile(
                contentPadding: EdgeInsets.zero,
                leading: Icon(Icons.pets),
                title: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Petifies Session',
                      style: TextStyle(
                        fontSize: 16,
                      ),
                    ),
                    statusWidget(status),
                  ],
                ),
                subtitle: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      "From " + formatter.format(fromTime),
                      style: TextStyle(
                        fontSize: 12,
                      ),
                    ),
                    Text(
                      "To " + formatter.format(toTime),
                      style: TextStyle(
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  if (showListProposalButton)
                    TextButton(
                      child: const Text(
                        'List proposals',
                        style: TextStyle(
                          fontSize: 12,
                        ),
                      ),
                      onPressed: () {
                        showModalBottomSheet(
                            context: context,
                            isScrollControlled: true,
                            useSafeArea: true,
                            barrierColor:
                                Theme.of(context).scaffoldBackgroundColor,
                            builder: (context) {
                              return ListProposalsForSessionScreen(
                                sessionId: id,
                              );
                            });
                      },
                      style: TextButton.styleFrom(
                          padding: EdgeInsets.only(
                        right: 12,
                      )),
                    ),
                  TextButton(
                    child: const Text(
                      'Make a proposal',
                      style: TextStyle(
                        fontSize: 12,
                      ),
                    ),
                    onPressed: () {
                      showModalBottomSheet(
                          context: context,
                          isScrollControlled: true,
                          useSafeArea: true,
                          barrierColor:
                              Theme.of(context).scaffoldBackgroundColor,
                          builder: (context) {
                            return PetifiesCreateProposalScreen(
                              sessionId: id,
                            );
                          });
                    },
                    style: ButtonStyle(
                      padding: MaterialStatePropertyAll(EdgeInsets.zero),
                    ),
                  ),
                ],
              )
            ],
          ),
        ),
      ),
    );
  }
}
