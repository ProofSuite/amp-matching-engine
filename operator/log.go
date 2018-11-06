package operator

// import (
// 	"github.com/Proofsuite/amp-matching-engine/types"
// 	"github.com/Proofsuite/amp-matching-engine/utils"
// 	logging "github.com/op/go-logging"
// )

// type OperatorLogger struct {
// 	*logging.Logger
// 	operatorLogger *logging.Logger
// }

// func NewOperatorLogger() *OperatorLogger {
// 	return &OperatorLogger{
// 		utils.StdoutLogger,
// 		utils.OperatorMessagesLogger,
// 	}
// }

// func (l *OperatorLogger) Log(msg string, m *types.Matches) {
// 	l.operatorLogger.Infof("%v: %v", msg, m.String())
// }

// func (l *OperatorLogger) LogTxError(m *types.Matches) {
// 	l.operatorLogger.Error("Transaction Failed: ", utils.JSON(m))
// }

// func (l *OperatorLogger) LogTxSuccess(m *types.Matches) {
// 	l.operatorLogger.Infof("Transaction Success: ", m.String())
// }

// func (l *OperatorLogger) LogMessageIn(msg *types.OperatorMessage) {
// 	l.operatorLogger.Infof("Receiving %v/%v message", msg.MessageType, msg.Matches.String())
// }

// func (l *OperatorLogger) LogMessageOut(msg *types.OperatorMessage) {
// 	l.operatorLogger.Infof("Sending %v/%v message", msg.MessageType, msg.Matches.String())
// }
