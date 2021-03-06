package p

import (
	"cloud.google.com/go/firestore"
	"context"
)

var newsletters *firestore.CollectionRef
var tldrs *firestore.CollectionRef

type Newsletter struct {
	NewsIds []string `firestore:"news_ids"`
	Test bool `firestore:"test"`
}

type TLDR struct {
	Summary string `firestore:"summary"`
	Title string `firestore:"title"`
	TopImage string `firestore:"top_image"`
	Url string `firestore:"url"`
}

func MaybeInitNewslettersCollection() error {
	ctx := context.Background()
	if newsletters == nil {
		firestoreClient, err := firestore.NewClient(ctx, mailSenderProjectId)
		if err != nil {
			return err
		}

		newsletters = firestoreClient.Collection("newsletters")
	}
	if tldrs == nil {
		firestoreClient, err := firestore.NewClient(ctx, tldrProjectId)
		if err != nil {
			return err
		}

		tldrs = firestoreClient.Collection("urls")
	}
	return nil
}

func GetNewsletterById(newsletterId string) (Newsletter, error) {
	ctx := context.Background()
	if err := MaybeInitNewslettersCollection(); err != nil {
		return Newsletter{}, err
	}
	newsletterDoc := newsletters.Doc(newsletterId)
	newsletterDocSnap, err := newsletterDoc.Get(ctx)
	if err != nil {
		return Newsletter{}, err
	}
	var newsletter Newsletter
	if err := newsletterDocSnap.DataTo(&newsletter); err != nil {
		return Newsletter{}, err
	}
	return newsletter, nil
}

func GetTldrById(tldrId string) (TLDR, error) {
	ctx := context.Background()
	if err := MaybeInitNewslettersCollection(); err != nil {
		return TLDR{}, err
	}
	tldrDoc := tldrs.Doc(tldrId)
	tldrDocSnap, err := tldrDoc.Get(ctx)
	if err != nil {
		return TLDR{}, err
	}
	var tldr TLDR
	if err := tldrDocSnap.DataTo(&tldr); err != nil {
		return TLDR{}, err
	}
	return tldr, nil
}

