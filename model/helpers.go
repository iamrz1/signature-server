package model

import cerror "signature-server/error"

func (b *VerifySignatureReq) Validate() error {
	err := cerror.ValidationError{}
	if b.Message == nil {
		err.Add("message", "required")
	}

	if b.Signature == "" {
		err.Add("signature", "required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (b *CreateTransactionReq) Validate() error {
	err := cerror.ValidationError{}
	if b.TxnData == "" {
		err.Add("txn", "required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}
