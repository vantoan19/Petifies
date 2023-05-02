import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/status/status.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';

class PetifiesProposal extends ConsumerWidget {
  const PetifiesProposal({super.key});

  Widget getStatusWidget(PetifiesProposalStatus status) {
    switch (status) {
      case PetifiesProposalStatus.ACCEPTED:
        return Status(
          color: Themes.greenColor,
          label: "ACCEPTED",
          textSize: 14,
          textColor: Themes.whiteColor,
        );
      case PetifiesProposalStatus.CANCELLED:
        return Status(
          color: Themes.redColor,
          label: "CANCELLED",
          textSize: 12,
          textColor: Themes.blackColor,
        );
      case PetifiesProposalStatus.REJECTED:
        return Status(
          color: Themes.redColor,
          label: "REJECTED",
          textSize: 12,
          textColor: Themes.blackColor,
        );
      case PetifiesProposalStatus.SESSION_CLOSED:
        return Status(
          color: Themes.greyColor,
          label: "SESSION CLOSED",
          textSize: 12,
          textColor: Themes.whiteColor,
        );
      case PetifiesProposalStatus.WAITING_FOR_ACCEPTANCE:
        return Status(
          color: Themes.blueColor,
          label: "WAITING FOR ACCEPTANCE",
          textSize: 12,
          textColor: Themes.whiteColor,
        );
      default:
        return Status(
          color: Themes.greyColor,
          label: "UNKNOWN",
          textSize: 12,
          textColor: Themes.whiteColor,
        );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final author = ref.read(proposalProvider.select((value) => value.author));
    final proposal =
        ref.read(proposalProvider.select((value) => value.proposal));
    final createdAt =
        ref.read(proposalProvider.select((value) => value.createdAt));
    final status = ref.read(proposalProvider.select((value) => value.status));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            UserAvatar(
              userAvatar: author.userAvatar,
              padding: EdgeInsets.only(right: 16),
            ),
            Column(
              mainAxisSize: MainAxisSize.max,
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  author.firstName + " " + author.lastName,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                ),
                Text(
                  StringUtils.stringifyTime(createdAt),
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w300,
                    color: Theme.of(context)
                        .colorScheme
                        .secondary
                        .withOpacity(0.8),
                  ),
                ),
              ],
            )
          ],
        ),
        Padding(
          padding: const EdgeInsets.only(top: 14.0),
          child: getStatusWidget(status),
        ),
        Padding(
          padding: const EdgeInsets.only(top: 14.0),
          child: Text(
            proposal,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w400,
              color: Theme.of(context).colorScheme.secondary,
              height: 1.35,
            ),
          ),
        ),
      ],
    );
  }
}
