while [[ 1 ]]; do
	m=`date +%M`
	am start -n com.ss.android.lark/com.ss.android.lark.main.app.MainActivity
	sleep 10
	am start -n com.termux/com.termux.app.TermuxActivity	
	sleep 10
done
