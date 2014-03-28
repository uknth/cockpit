<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Auth extends CI_Controller {

	function index(){
		redirectIfNoSession($this);
		$this->layout->setLayout('layouts/main');
		$this->layout->view('home/homepage',array(
			'js_file' => 'AdminLTE/dashboard'
		));
	}

	public function login(){
		//echo '<pre>'; print_r(apache_get_modules());echo '</pre>';exit;
		redirectIfSession($this);
		$this->layout->setLayout('layouts/auth');
		$this->layout->view('auth/login');
	}

	public function register(){
		redirectIfSession($this);
		$this->layout->setLayout('layouts/auth');
		$this->layout->view('auth/register');
	}

	public function doRegister(){
		redirectIfSession($this);

		$this->load->library('form_validation');

		$this->form_validation->set_rules('password', 'Password', 'required');
		$this->form_validation->set_rules('email', 'Email', 'required|valid_email');

		if ($this->form_validation->run() == FALSE){
			echo json_encode(array(
				'status' => 'error'
				,'msg' => $this->form_validation->error_string()
			));
			exit;
		}

		$post = $this->input->post();

		//lets load user model
		//u is alias of User_db
		$this->load->model('User_db','u');

		//lets create new user by inserting into db
		$user_id = $this->u->insertNewUser($post['email'],$post['password']);

		if($user_id <= 0){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> $user_id == 0 ? 'Somer error occured. Please try after some time.' : 'An user is already existing with this account. Please user correct password to login.'
			));
			exit;
		}

		$random_string = substr(sha1(rand()), 0, 15);

		//lets set session
		$this->session->set_userdata(array(
			'user_id' => $user_id
			,'email' => $post['email']
		));

		//LETS REGISTER IN GO SERVER
		$status = post_to_go('auth','',array('USER_ID' => $user_id, 'TOKEN' => $random_string),array('USER_ID : '.$user_id, 'TOKEN : '.$post['email']));

		echo json_encode(array(
			'status' => 'success'
			,'message' => 'You are successfully registered.'
			,'user_id' => $user_id
			,'token' => $random_string
		));
		exit;
	}

	function doLogin(){
		redirectIfSession($this);

		$this->load->library('form_validation');

		$this->form_validation->set_rules('password', 'Password', 'required');
		$this->form_validation->set_rules('email', 'Email', 'required|valid_email');

		if ($this->form_validation->run() == FALSE){
			echo json_encode(array(
				'status' => 'error'
				,'msg' => $this->form_validation->error_string()
			));
			exit;
		}

		$post = $this->input->post();

		//lets load user model
		//u is alias of User_db
		$this->load->model('User_db','u');

		//lets get user details from db
		//$data = $this->u->getUserDetails($post['email']);
		$data = (object) array("user_id"=>3,"password"=>"anilanl");
		if(!is_object($data)){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> 'This email is not registered with us. Please register'
			));
			exit;
		}

		if($data->password != $post['password']){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> 'Please provide correct password'
			));
			exit;
		}

		$random_string = substr(sha1(rand()), 0, 15);

		//lets set session
		$this->session->set_userdata(array(
			'user_id' => $data->user_id
			,'email' => $post['email']
			,'token' => $random_string
		));

		//LETS REGISTER IN GO SERVER
		$status = post_to_go('auth','',array('USER_ID' => $data->user_id, 'TOKEN' => $post['email']),array('USER_ID: '.$data->user_id, 'TOKEN: '.$post['email']));

		echo json_encode(array(
			'status' => 'success'
			,'message' => 'You are successfully logged-in.'
			,'user_id' => $user_id
			,'token' => $random_string
			,'result' => $status
		));
		exit;
	}

	function logout(){
		$this->session->destroy();
		redirect(base_url());
	}

	function test(){
		// Replace with real server API key from Google APIs          
        $apiKey = "AIzaSyCNJeYio3_-9tRF8fjuUep7BeV0c0rE8TU";
        
        // Replace with real client registration IDs
        $registrationIDs = array( "APA91bGvKhOyHmvldS2yegp3gJ8swrnvGLTGHmybI7CK4Ui5o-A2un2o-ztegeuwX5zf3W1R96FU8E8t4rMX2DooTmAblc_yjCBYUsuQ5AIrNO-_gnxBQkSwPE0NiTEn5OpLmzDHo1xpjeL3iVXhePuJZaeFoAZs3-wRTCqgc9iq_UDioTrDv0E");
        
        $url = 'https://android.googleapis.com/gcm/send';
		$serverApiKey = $apiKey;
		$reg = "APA91bHaslUVsUrVDsUpOZ-FJ7AMUaWXnh7kfG4ztpwVVbztBSZcu55kG24F4xMl4qEBASZlCyYGsyP_f9-nsIAwbw-Lhb1bjZa1FKIzSA7-4cw39w2gwg00uyywM7xJSupx6AUG8EX-3fbiWtyFPoCv5A1IK3XIVG6cQKmoRm_roBfK1On-Chk"; // registration id

		$headers = array(
		'Content-Type:application/json',
		'Authorization:key=' . $serverApiKey
		);

		$data = array(
		'registration_ids' => array($reg),
		'data' => array(
		'message' => 'Hello, World!'
		));

		print (json_encode($data) . "\n\n");

		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, $url);
		if ($headers)
		curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
		curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
		curl_setopt($ch, CURLOPT_POST, true);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));

		$response = curl_exec($ch);

		$header = curl_getinfo ( $ch );


		curl_close($ch);

		print_r ($response);

		echo '<pre>'; print_r($header);echo '</pre>';
	}
}